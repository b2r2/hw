package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

const (
	ExitError = 1
	ExitOk    = 0
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return ExitError
	}
	for k, v := range env {
		if err := os.Unsetenv(k); err != nil {
			log.Println(err)
			return ExitError
		}
		if !v.NeedRemove {
			if err := os.Setenv(k, v.Value); err != nil {
				log.Println(err)
				return ExitError
			}
		}
	}
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Env = os.Environ()

	if err := c.Run(); err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			return e.ExitCode()
		}
		return ExitError
	}
	return ExitOk
}
