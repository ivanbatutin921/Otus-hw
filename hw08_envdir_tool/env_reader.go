package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := dir + "/" + file.Name()
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		line, err := reader.ReadString('\n')
		if !errors.Is(err, io.EOF) && err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, " \t")
		if len(line) == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
		} else {
			env[file.Name()] = EnvValue{Value: line}
		}
	}
	return env, nil
}
