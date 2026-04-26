package executors

import (
	"fmt"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

func runToolByID(
	context runner.Context,
	workdir string,
	toolID string,
	arguments ...string,
) (output string, err error) {
	tool, found := context.Effective.ToolByID(toolID)
	if !found {
		return "", fmt.Errorf("unknown tool %q", toolID)
	}

	capability, found := context.ToolCapabilities[toolID]
	if !found {
		return "", fmt.Errorf("unknown tool capability %q", toolID)
	}

	return runtime.RunToolCommand(workdir, context.GoEnvironment, tool, capability, arguments...)
}

func runCommandOutput(
	workdir string,
	environment map[string]string,
	name string,
	arguments ...string,
) (output string, err error) {
	result, err := runtime.RunCommand(runtime.CommandRequest{
		Directory:   workdir,
		Environment: environment,
		Name:        name,
		Arguments:   append([]string{}, arguments...),
	})
	return runtime.CommandOutput(result, err)
}
