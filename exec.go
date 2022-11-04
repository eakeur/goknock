package main

import (
	"fmt"
	"os/exec"
	"strings"
)

const ErrPlaceholder = `
Command error status: 
%w
		
Command error output:
%s`

type Executer func(command string) (string, error)

func executer() Executer {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}

	return func(command string) (string, error) {
		cmd := exec.Command("powershell", command)
		cmd.Stderr = stderr
		cmd.Stdout = stdout

		defer func(stdout, stderr *strings.Builder) {
			stdout.Reset()
			stderr.Reset()
		}(stdout, stderr)

		err := cmd.Start()
		if err != nil {
			return "", fmt.Errorf(ErrPlaceholder, err, stderr.String())
		}

		err = cmd.Wait()
		if err != nil {
			return "", fmt.Errorf(ErrPlaceholder, err, stderr.String())
		}

		return stdout.String(), nil
	}
}
