package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func executeCommand(client *ssh.Client, command string, failIfNonzero bool) (string, int64, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", -1, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	outputBytes, err := session.CombinedOutput(command)
	outputStr := string(outputBytes)

	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode := int64(exitErr.ExitStatus())
			if failIfNonzero && exitCode != 0 {
				return outputStr, exitCode, fmt.Errorf("command exited with non-zero status: %d\nOutput: %s",
					exitCode, outputStr)
			}
			return outputStr, exitCode, nil
		}
		return outputStr, -1, fmt.Errorf("failed to execute command: %w", err)
	}

	return outputStr, 0, nil
}
