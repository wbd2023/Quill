package commandrun

import (
	"fmt"
	"time"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
)

// ToolByID runs a tool identified by toolID and returns its result.
func ToolByID(
	context runner.Context,
	workDir string,
	toolID string,
	arguments ...string,
) (result runtime.CommandResult, err error) {
	tool, found := context.Tools[toolID]
	if !found {
		return runtime.CommandResult{}, fmt.Errorf("unknown tool %q", toolID)
	}

	return runtime.RunCommand(runtime.CommandRequest{
		Name:             tool.Command,
		Arguments:        append([]string{}, arguments...),
		Environment:      context.GoEnvironment,
		Directory:        workDir,
		Timeout:          time.Duration(tool.TimeoutSeconds) * time.Second,
		OutputLimitBytes: tool.OutputLimitBytes,
	})
}

// Output runs a command and returns its result.
func Output(
	workDir string,
	environment map[string]string,
	name string,
	arguments ...string,
) (result runtime.CommandResult, err error) {
	return runtime.RunCommand(runtime.CommandRequest{
		Name:        name,
		Arguments:   append([]string{}, arguments...),
		Environment: environment,
		Directory:   workDir,
	})
}

// BuildStyleResult projects a runtime.CommandResult onto the style.CommandResult the
// report layer consumes (drops Output, which the report does not need).
func BuildStyleResult(result runtime.CommandResult) (commandResult style.CommandResult) {
	return style.CommandResult{
		ExitCode:  result.ExitCode,
		TimedOut:  result.TimedOut,
		Truncated: result.Truncated,
	}
}
