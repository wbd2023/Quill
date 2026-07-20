package commandrun

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/process"
)

// ToolByID runs a tool identified by toolID and returns its result.
func ToolByID(
	ctx context.Context,
	context execution.RunContext,
	workDir string,
	toolID string,
	arguments ...string,
) (result process.CommandResult, err error) {
	tool, found := context.Tools[toolID]
	if !found {
		return process.CommandResult{}, fmt.Errorf("unknown tool %q", toolID)
	}

	return process.RunCommand(ctx, process.CommandRequest{
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
	ctx context.Context,
	workDir string,
	environment map[string]string,
	name string,
	arguments ...string,
) (result process.CommandResult, err error) {
	return process.RunCommand(ctx, process.CommandRequest{
		Name:        name,
		Arguments:   slices.Clone(arguments),
		Environment: environment,
		Directory:   workDir,
	})
}
