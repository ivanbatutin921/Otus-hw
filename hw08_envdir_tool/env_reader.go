package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
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
	if dir == "" {
		return nil, ErrEmptyDirectoryPath
	}

	environment := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()

		if err = checkFileName(fileName); err != nil {
			return nil, err
		}

		filePath := dir + "/" + fileName

		fileParams, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		if fileParams.Size() == 0 {
			envValue := EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			environment[fileName] = envValue
			continue
		}

		openedFile, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		environmentValue, err := firstLine(openedFile)
		if !errors.Is(err, io.EOF) && err != nil {
			return nil, err
		}

		envValue := EnvValue{
			Value:      environmentValue,
			NeedRemove: false,
		}
		environment[fileName] = envValue
	}
	return environment, nil
}

func firstLine(file *os.File) (string, error) {
	reader := bufio.NewReader(file)

	firstString, err := reader.ReadString('\n')
	if err != io.EOF && err != nil {
		return "", err
	}

	environmentValue := strings.TrimRight(strings.ReplaceAll(firstString, "\000", "\n"), "\n")

	return environmentValue, err
}

func checkFileName(name string) error {
	checkFilename, err := regexp.MatchString("=", name)
	if checkFilename {
		return ErrWrongFileName
	}
	return err
}
