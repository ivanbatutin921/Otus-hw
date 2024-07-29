package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = os.Environ()
	for key, value := range env {
		if value.NeedRemove {
			command.Env = removeEnv(command.Env, key)
		} else {
			command.Env = append(command.Env, key+"="+value.Value)
		}
	}
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		returnCode = 1
	}
	return
}

func removeEnv(env []string, key string) []string {
	for i, v := range env {
		if strings.HasPrefix(v, key+"=") {
			return append(env[:i], env[i+1:]...)
		}
	}
	return env
}
