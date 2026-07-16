package runtime

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// Default limits applied when a request omits its own timeout or output limit.
const (
	defaultCommandTimeout   = 120 * time.Second
	defaultOutputLimitBytes = 1 << 20
)

const anyExecutePermission os.FileMode = 0o111

/* -------------------------------------------- Types ------------------------------------------- */

// CommandRequest represents a command to execute: the executable name, arguments, environment,
// working directory, and execution limits.
type CommandRequest struct {
	Name      string
	Arguments []string

	Environment map[string]string
	Directory   string

	Timeout          time.Duration
	OutputLimitBytes int64
}

// CommandResult represents the outcome of running a command: captured output, exit status, timeout
// status, and output truncation.
type CommandResult struct {
	Output   string
	ExitCode int

	TimedOut  bool
	Truncated bool
}

/* -------------------------------------- Command Execution ------------------------------------- */

// RunCommand executes the command described by request, applying its timeout and output limit, and
// returns the captured output and exit status. A non-zero exit or timeout is returned as a
// CommandError wrapping the underlying cause.
func RunCommand(request CommandRequest) (result CommandResult, err error) {
	path, err := ResolveCommandPath(request.Environment, request.Name)
	if err != nil {
		return CommandResult{}, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		resolveTimeout(request.Timeout),
	)
	defer cancel()

	command := exec.CommandContext(ctx, path, request.Arguments...)
	command.Dir = resolveDirectory(request.Directory)
	command.Env = mergeEnvironment(request.Environment)

	buffer := &limitedBuffer{limit: resolveOutputLimit(request.OutputLimitBytes)}
	command.Stdout = buffer
	command.Stderr = buffer

	err = command.Run()
	result = CommandResult{
		Output:    buffer.String(),
		ExitCode:  resolveExitCode(err),
		TimedOut:  ctx.Err() != nil,
		Truncated: buffer.truncated,
	}

	if err == nil {
		return result, nil
	}

	if result.TimedOut {
		err = fmt.Errorf("%w: %w", context.DeadlineExceeded, err)
	}

	return result, CommandError{
		Name:      request.Name,
		Arguments: slices.Clone(request.Arguments),
		Err:       err,
	}
}

/* --------------------------------------- Path Resolution -------------------------------------- */

// ResolveCommandPath resolves command to an executable path. It honours the provided environment's
// PATH when set, otherwise falls back to exec.LookPath. Absolute paths and paths containing a
// separator are returned as-is.
func ResolveCommandPath(
	environment map[string]string,
	command string,
) (resolved string, err error) {
	path, ok := environment["PATH"]
	if !ok {
		path = os.Getenv("PATH")
	}
	if path == "" {
		return exec.LookPath(command)
	}

	if filepath.IsAbs(command) || strings.ContainsRune(command, os.PathSeparator) {
		return command, nil
	}

	for _, directory := range filepath.SplitList(path) {
		candidate := filepath.Join(directory, command)
		if isExecutableFile(candidate) {
			return candidate, nil
		}
	}

	return "", exec.ErrNotFound
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func mergeEnvironment(extra map[string]string) (environment []string) {
	environment = os.Environ()
	for _, key := range slices.Sorted(maps.Keys(extra)) {
		environment = append(environment, key+"="+extra[key])
	}

	return environment
}

func resolveDirectory(directory string) (resolved string) {
	if directory == "" {
		return "."
	}

	return directory
}

func resolveTimeout(timeout time.Duration) (duration time.Duration) {
	if timeout <= 0 {
		return defaultCommandTimeout
	}

	return timeout
}

func resolveOutputLimit(limit int64) (resolved int64) {
	if limit <= 0 {
		return defaultOutputLimitBytes
	}

	return limit
}

func resolveExitCode(err error) (code int) {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return -1
}

func isExecutableFile(path string) (found bool) {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode()&anyExecutePermission != 0
}
