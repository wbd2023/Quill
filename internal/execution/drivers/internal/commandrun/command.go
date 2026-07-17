package commandrun

import (
	"fmt"
	"slices"
	"time"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/process"
)

// ToolByID runs a tool identified by toolID and returns its result.
func ToolByID(
	context execution.RunContext,
	workDir string,
	toolID string,
	arguments ...string,
) (result process.CommandResult, err error) {
	tool, found := context.Tools[toolID]
	if !found {
		return process.CommandResult{}, fmt.Errorf("unknown tool %q", toolID)
	}

	return process.RunCommand(process.CommandRequest{
		Name:             tool.Command,
		Arguments:        slices.Clone(arguments),
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
) (result process.CommandResult, err error) {
	return process.RunCommand(process.CommandRequest{
		Name:        name,
		Arguments:   slices.Clone(arguments),
		Environment: environment,
		Directory:   workDir,
	})
}
