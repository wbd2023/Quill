package runtime

import (
	"fmt"
	"strings"
)

type CommandError struct {
	Name      string
	Arguments []string
	Result    CommandResult
	Err       error
}

func (err CommandError) Error() (message string) {
	switch {
	case err.Result.TimedOut:
		return fmt.Sprintf("%s timed out", commandText(err.Name, err.Arguments))
	case err.Err != nil:
		return fmt.Sprintf("%s failed: %v", commandText(err.Name, err.Arguments), err.Err)
	default:
		return fmt.Sprintf("%s failed with exit code %d",
			commandText(err.Name, err.Arguments),
			err.Result.ExitCode,
		)
	}
}

func (err CommandError) Unwrap() (wrapped error) {
	return err.Err
}

func commandText(name string, arguments []string) (text string) {
	if len(arguments) == 0 {
		return name
	}

	return name + " " + strings.Join(arguments, " ")
}
