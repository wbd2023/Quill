package process

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// CommandError represents a failed command execution: the command name, its arguments, and the
// underlying error.
type CommandError struct {
	Name      string
	Arguments []string
	Err       error
}

func (err CommandError) Error() (message string) {
	command := err.Name
	if len(err.Arguments) > 0 {
		command += " " + strings.Join(err.Arguments, " ")
	}

	if errors.Is(err.Err, context.DeadlineExceeded) {
		return fmt.Sprintf("%s timed out", command)
	}

	return fmt.Sprintf("%s failed: %v", command, err.Err)
}

func (err CommandError) Unwrap() (wrapped error) {
	return err.Err
}
