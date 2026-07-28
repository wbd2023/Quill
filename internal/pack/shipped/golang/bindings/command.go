package bindings

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/process"
)

// runTool runs a pinned Go tool identified by toolID in workDir against the repository's Go
// environment, applying the tool's configured timeout and output limit. It is the single place
// that turns a tool ID plus the run context into a correctly configured process call, so each Go
// target closure reuses one setup path instead of reconstructing it.
func runTool(
	ctx context.Context,
	run execution.RunContext,
	workDir string,
	toolID string,
	arguments ...string,
) (result process.CommandResult, err error) {
	tool, found := run.Tools[toolID]
	if !found {
		return process.CommandResult{}, fmt.Errorf("unknown tool %q", toolID)
	}

	return process.RunCommand(ctx, process.CommandRequest{
		Name:             tool.Command,
		Arguments:        slices.Clone(arguments),
		Environment:      process.EnvironmentInherit,
		Variables:        run.GoEnvironment,
		Directory:        workDir,
		Timeout:          time.Duration(tool.TimeoutSeconds) * time.Second,
		OutputLimitBytes: tool.OutputLimitBytes,
	})
}

// runGo runs a raw Go toolchain command (go, gofmt) in workDir against the repository's Go
// environment. Unlike runTool it has no pinned tool entry, so it mirrors direct toolchain
// invocation with no timeout or output limit.
func runGo(
	ctx context.Context,
	run execution.RunContext,
	workDir string,
	name string,
	arguments ...string,
) (result process.CommandResult, err error) {
	return process.RunCommand(ctx, process.CommandRequest{
		Name:        name,
		Arguments:   slices.Clone(arguments),
		Environment: process.EnvironmentInherit,
		Variables:   run.GoEnvironment,
		Directory:   workDir,
	})
}
