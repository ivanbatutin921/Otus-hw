package main

import (
	"log"
	"os"
)

func main() {
	cmdArgs := os.Args
	dirPath := cmdArgs[1]

	env, err := ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	resultCode := RunCmd(cmdArgs, env)
	os.Exit(resultCode)
}
