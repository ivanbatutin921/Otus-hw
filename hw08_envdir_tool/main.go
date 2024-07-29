package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir <dir> <command> [args...]")
		return
	}
	dir := os.Args[1]
	cmd := os.Args[2:]
	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	returnCode := RunCmd(cmd, env)
	os.Exit(returnCode)
}
