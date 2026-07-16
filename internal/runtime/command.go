package runtime

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"slices"
	"time"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// Default limits applied when a request omits its own timeout or output limit.
const (
	defaultCommandTimeoutSeconds = 120
	defaultOutputLimitBytes      = 1 << 20
)

/* -------------------------------------------- Types ------------------------------------------- */

// CommandRequest represents a command to execute: the executable name, arguments, environment,
// working directory, and execution limits.
type CommandRequest struct {
	Name             string
	Arguments        []string
	Environment      map[string]string
	Directory        string
	TimeoutSeconds   int
	OutputLimitBytes int64
}

// CommandResult represents the outcome of running a command: captured output, exit status, timeout
// status, and output truncation.
type CommandResult struct {
	Output    string
	ExitCode  int
	TimedOut  bool
	Truncated bool
}

/* -------------------------------------- Command Execution ------------------------------------- */

// RunCommand executes the command described by request, applying its timeout and output limit, and
// returns the captured output and exit status. A non-zero exit or timeout is returned as a
// CommandError carrying the result.
func RunCommand(request CommandRequest) (result CommandResult, err error) {
	path, err := ResolveCommandPath(request.Name, request.Environment)
	if err != nil {
		return CommandResult{}, err
	}

	seconds := request.TimeoutSeconds
	if seconds <= 0 {
		seconds = defaultCommandTimeoutSeconds
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	command := exec.CommandContext(ctx, path, request.Arguments...)
	command.Dir = request.Directory
	if command.Dir == "" {
		command.Dir = "."
	}
	command.Env = append(os.Environ(), environmentEntries(request.Environment)...)

	limit := request.OutputLimitBytes
	if limit <= 0 {
		limit = defaultOutputLimitBytes
	}
	buffer := &limitedBuffer{limit: limit}
	command.Stdout = buffer
	command.Stderr = buffer

	err = command.Run()
	result = CommandResult{
		Output:    buffer.String(),
		ExitCode:  exitCode(err),
		TimedOut:  ctx.Err() != nil,
		Truncated: buffer.truncated,
	}
	if err != nil {
		return result, CommandError{
			Name:      request.Name,
			Arguments: append([]string{}, request.Arguments...),
			Result:    result,
			Err:       err,
		}
	}

	return result, nil
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func environmentEntries(environment map[string]string) (entries []string) {
	if len(environment) == 0 {
		return nil
	}

	keys := make([]string, 0, len(environment))
	for key := range environment {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	entries = make([]string, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, key+"="+environment[key])
	}

	return entries
}

func exitCode(err error) (code int) {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return -1
}
