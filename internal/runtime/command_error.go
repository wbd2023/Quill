package runtime

import (
	"fmt"
	"strings"
)

// CommandError represents a failed command execution: the command name, its arguments, whether it
// timed out, and the underlying error.
type CommandError struct {
	Name      string
	Arguments []string
	TimedOut  bool
	Err       error
}

func (err CommandError) Error() (message string) {
	formatted := formatCommand(err.Name, err.Arguments)

	if err.TimedOut {
		return fmt.Sprintf("%s timed out", formatted)
	}

	return fmt.Sprintf("%s failed: %v", formatted, err.Err)
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
