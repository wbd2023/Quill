package commandrun

import (
	"fmt"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
)

// ToolByID runs a tool's command and returns its stdout.
func ToolByID(
	context runner.Context,
	workDir string,
	toolID string,
	arguments ...string,
) (output string, err error) {
	tool, found := context.Tools[toolID]
	if !found {
		return "", fmt.Errorf("unknown tool %q", toolID)
	}

	result, err := runtime.RunCommand(runtime.CommandRequest{
		Directory:        workDir,
		Environment:      context.GoEnvironment,
		Name:             tool.Command,
		Arguments:        append([]string{}, arguments...),
		TimeoutSeconds:   tool.TimeoutSeconds,
		OutputLimitBytes: tool.OutputLimitBytes,
	})
	return result.Output, err
}

// Output runs a command and returns its stdout.
func Output(
	workDir string,
	environment map[string]string,
	name string,
	arguments ...string,
) (output string, err error) {
	result, err := runtime.RunCommand(runtime.CommandRequest{
		Directory:   workDir,
		Environment: environment,
		Name:        name,
		Arguments:   append([]string{}, arguments...),
	})
	return result.Output, err
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
