package process

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// CommandError represents a failed command execution: the command name, its arguments, and the
// underlying error. The underlying error distinguishes parent cancellation
// (context.Canceled), child timeout (context.DeadlineExceeded), and ordinary exit failure, so the
// two termination causes are never collapsed into a generic message.
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

	switch {
	case errors.Is(err.Err, context.Canceled):
		return fmt.Sprintf("%s canceled", command)
	case errors.Is(err.Err, context.DeadlineExceeded):
		return fmt.Sprintf("%s timed out", command)
	default:
		return fmt.Sprintf("%s failed: %v", command, err.Err)
	}
}

func (err CommandError) Unwrap() (wrapped error) {
	return err.Err
}
