package executors

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

func fileCommandExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	detail, found := spec.FileCommandExecution()
	if !found {
		return contract.ExecutionResult{}, errors.New("file-command executor received empty spec")
	}

	files, err := runner.CollectFileSetFiles(context, detail.FileSet)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	if len(files) == 0 {
		return contract.ExecutionResult{}, nil
	}

	tool, found := context.Effective.ToolByID(detail.ToolID)
	if !found {
		return contract.ExecutionResult{}, errUnknownTool(detail.ToolID)
	}

	capability, found := context.ToolCapabilities[detail.ToolID]
	if !found {
		return contract.ExecutionResult{}, errUnknownTool(detail.ToolID)
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
