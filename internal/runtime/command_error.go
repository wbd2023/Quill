package runtime

import (
	"fmt"
	"strings"
)

// CommandError represents a failed command execution: the command name, its arguments, the captured
// result, and the underlying error.
type CommandError struct {
	Name      string
	Arguments []string
	Result    CommandResult
	Err       error
}

func (err CommandError) Error() (message string) {
	switch {
	case err.Result.TimedOut:
		return fmt.Sprintf("%s timed out", formatCommand(err.Name, err.Arguments))

	case err.Err != nil:
		return fmt.Sprintf("%s failed: %v", formatCommand(err.Name, err.Arguments), err.Err)

	default:
		return fmt.Sprintf("%s failed with exit code %d",
			formatCommand(err.Name, err.Arguments),
			err.Result.ExitCode,
		)
	}
}

func (err CommandError) Unwrap() (wrapped error) {
	return err.Err
}

func formatCommand(name string, arguments []string) (command string) {
	if len(arguments) == 0 {
		return name
	}

	return name + " " + strings.Join(arguments, " ")
}
