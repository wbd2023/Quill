package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	defaultCommandTimeoutSeconds = 120
	defaultOutputLimitBytes      = 1 << 20
)

/* -------------------------------------------- Types ------------------------------------------- */

type CommandRequest struct {
	Directory        string
	Environment      map[string]string
	Name             string
	Arguments        []string
	TimeoutSeconds   int
	OutputLimitBytes int64
}

type CommandResult struct {
	Output    string
	ExitCode  int
	TimedOut  bool
	Truncated bool
}

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

type limitedBuffer struct {
	builder   strings.Builder
	limit     int64
	written   int64
	truncated bool
}

/* -------------------------------------- Command Execution ------------------------------------- */

func RunCommand(request CommandRequest) (result CommandResult, err error) {
	commandPath, err := lookupCommandPath(request.Name, request.Environment)
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

func RunToolCommand(
	directory string,
	environment map[string]string,
	tool contract.Tool,
	capability toolchain.Capability,
	arguments ...string,
) (output string, err error) {
	result, err := RunToolCommandResult(
		directory,
		environment,
		tool,
		capability,
		arguments...,
	)
	return CommandOutput(result, err)
}

func RunToolCommandResult(
	directory string,
	environment map[string]string,
	tool contract.Tool,
	capability toolchain.Capability,
	arguments ...string,
) (result CommandResult, err error) {
	return RunCommand(CommandRequest{
		Directory:        directory,
		Environment:      environment,
		Name:             capability.Command,
		Arguments:        append([]string{}, arguments...),
		TimeoutSeconds:   tool.TimeoutSeconds,
		OutputLimitBytes: tool.OutputLimitBytes,
	})
}

func ContractCommandResult(result CommandResult) (contractResult contract.CommandResult) {
	return contract.CommandResult{
		ExitCode:  result.ExitCode,
		TimedOut:  result.TimedOut,
		Truncated: result.Truncated,
	}
}

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

func commandText(name string, arguments []string) (text string) {
	if len(arguments) == 0 {
		return name
	}

	return name + " " + strings.Join(arguments, " ")
}

func (buffer *limitedBuffer) Write(data []byte) (count int, err error) {
	count = len(data)
	remaining := buffer.limit - buffer.written
	if remaining <= 0 {
		buffer.truncated = true
		return count, nil
	}

	if int64(len(data)) > remaining {
		data = data[:int(remaining)]
		buffer.truncated = true
	}

	buffer.written += int64(len(data))
	_, _ = buffer.builder.Write(data)
	return count, nil
}

func (buffer *limitedBuffer) String() (output string) {
	return buffer.builder.String()
}

func environmentEntries(environment map[string]string) (values []string) {
	if len(environment) == 0 {
		return nil
	}

	keys := make([]string, 0, len(environment))
	for key := range environment {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values = make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, key+"="+environment[key])
	}

	return values
}
