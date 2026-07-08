package runtime

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"

	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// command constants.
const (
	defaultCommandTimeoutSeconds = 120
	defaultOutputLimitBytes      = 1 << 20
)

/* -------------------------------------------- Types ------------------------------------------- */

// CommandRequest is command request.
type CommandRequest struct {
	Directory        string
	Environment      map[string]string
	Name             string
	Arguments        []string
	TimeoutSeconds   int
	OutputLimitBytes int64
}

// CommandResult is command result.
type CommandResult struct {
	Output    string
	ExitCode  int
	TimedOut  bool
	Truncated bool
}

/* -------------------------------------- Command Execution ------------------------------------- */

// RunCommand run command.
func RunCommand(request CommandRequest) (result CommandResult, err error) {
	commandPath, err := toolchain.ResolveCommandPath(request.Name, request.Environment)
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

// CommandOutput command output.
func CommandOutput(result CommandResult, err error) (output string, commandErr error) {
	if err == nil {
		return result.Output, nil
	}

	return result.Output, err
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
