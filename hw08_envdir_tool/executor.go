package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, value := range env {
		if value.NeedRemove {
			if err := os.Unsetenv(envName); err != nil {
				log.Fatal(err)
			}
		}
		if err := os.Setenv(envName, value.Value); err != nil {
			log.Fatal(err)
		}
	}

	command := exec.Command(cmd[2], cmd[3], cmd[4], cmd[5]) //nolint:gosec

	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		defaultErrCode := 100
		return defaultErrCode
	}
	return 0
}
