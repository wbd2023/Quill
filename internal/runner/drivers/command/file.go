package command

import (
	"errors"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func fileCommandDriver(
	context runner.Context,
	spec style.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result style.ExecutionResult, err error) {
	execution, found := spec.FileCommandExecution()
	if !found {
		return style.ExecutionResult{}, errors.New("file-command driver received empty spec")
	}

	files, err := runner.CollectFileSetFiles(context, execution.FileSet)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	if len(files) == 0 {
		return style.ExecutionResult{}, nil
	}

	tool, found := context.Effective.ToolByID(execution.ToolID)
	if !found {
		return style.ExecutionResult{}, errUnknownTool(execution.ToolID)
	}

	capability, found := context.ToolCapabilities[execution.ToolID]
	if !found {
		return style.ExecutionResult{}, errUnknownTool(execution.ToolID)
	}

	arguments := runner.FileCommandArguments(context.RepoRoot, spec)
	arguments = append(arguments, files...)
	commandResult, err := runtime.RunToolCommandResult(
		context.RepoRoot,
		context.ToolEnvironment,
		tool,
		capability,
		arguments...,
	)
	output, err := runtime.CommandOutput(commandResult, err)
	return style.ExecutionResult{
		Output:  output,
		Command: runtime.BuildStyleCommandResult(commandResult),
	}, err
}

func errUnknownTool(toolID string) (err error) {
	return errors.New("unknown tool " + toolID)
}
