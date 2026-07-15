package runtime

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// Default limits applied when a request omits its own timeout or output cap.
const (
	defaultCommandTimeoutSeconds = 120
	defaultOutputLimitBytes      = 1 << 20
)

/* -------------------------------------------- Types ------------------------------------------- */

// CommandRequest describes a command to run: its working directory, environment,
// executable, arguments, and execution limits.
type CommandRequest struct {
	Directory        string
	Environment      map[string]string
	Name             string
	Arguments        []string
	TimeoutSeconds   int
	OutputLimitBytes int64
}

// CommandResult holds the outcome of a command: its captured output, exit status, and
// whether it timed out or had its output truncated.
type CommandResult struct {
	Output    string
	ExitCode  int
	TimedOut  bool
	Truncated bool
}

/* -------------------------------------- Command Execution ------------------------------------- */

// RunCommand executes the command described by request, applying its timeout and output
// limit, and returns the captured output and exit status. A non-zero exit or timeout is
// returned as a CommandError carrying the result.
func RunCommand(request CommandRequest) (result CommandResult, err error) {
	commandPath, err := ResolveCommandPath(request.Name, request.Environment)
	if err != nil {
		return CommandResult{}, err
	}

	timeout := commandTimeout(request.TimeoutSeconds)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	command := exec.CommandContext(ctx, commandPath, request.Arguments...)
	command.Dir = commandDirectory(request.Directory)
	command.Env = append(os.Environ(), environmentEntries(request.Environment)...)

	buffer := &limitedBuffer{limit: commandOutputLimit(request.OutputLimitBytes)}
	command.Stdout = buffer
	command.Stderr = buffer

	runErr := command.Run()
	result = CommandResult{
		Output:    buffer.String(),
		ExitCode:  commandExitCode(runErr),
		TimedOut:  ctx.Err() != nil,
		Truncated: buffer.truncated,
	}
	if runErr != nil {
		return result, CommandError{
			Name:      request.Name,
			Arguments: append([]string{}, request.Arguments...),
			Result:    result,
			Err:       runErr,
		}
	}

	return result, nil
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func commandDirectory(directory string) (normalised string) {
	if directory == "" {
		return "."
	}

	return directory
}

func commandTimeout(seconds int) (timeout time.Duration) {
	if seconds <= 0 {
		seconds = defaultCommandTimeoutSeconds
	}

	return time.Duration(seconds) * time.Second
}

func commandOutputLimit(limit int64) (normalised int64) {
	if limit <= 0 {
		return defaultOutputLimitBytes
	}

	return limit
}

func commandExitCode(err error) (exitCode int) {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return -1
}
