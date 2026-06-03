package command

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

func fileCommandDriver(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.FileCommandExecution()
	if !found {
		return contract.ExecutionResult{}, errors.New("file-command driver received empty spec")
	}

	files, err := runner.CollectFileSetFiles(context, execution.FileSet)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	if len(files) == 0 {
		return contract.ExecutionResult{}, nil
	}

	tool, found := context.Effective.ToolByID(execution.ToolID)
	if !found {
		return contract.ExecutionResult{}, errUnknownTool(execution.ToolID)
	}

	capability, found := context.ToolCapabilities[execution.ToolID]
	if !found {
		return contract.ExecutionResult{}, errUnknownTool(execution.ToolID)
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
	return contract.ExecutionResult{
		Output:  output,
		Command: runtime.ContractCommandResult(commandResult),
	}, err
}

func errUnknownTool(toolID string) (err error) {
	return errors.New("unknown tool " + toolID)
}
