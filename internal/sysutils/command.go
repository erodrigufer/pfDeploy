package sysutils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ShCmd, run a shell command. args[0] is the command name, args[1:] are the
// command arguments.
func ShCmd(args ...string) (string, error) {
	// Command name.
	baseCmd := args[0]
	// Command arguments.
	cmdArgs := args[1:]

	cmd := exec.Command(baseCmd, cmdArgs...)
	// Buffers used to store the stdout and stderr output of command.
	var outb, errb bytes.Buffer
	// Redirect stdout and stderr to the byte buffers.
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("stderr: %s: error executing command: %w", errb.String(), err)
	}
	// Return the cmd stdout, if successful.
	return outb.String(), nil
}
