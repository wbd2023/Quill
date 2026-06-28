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
	commandResult, runErr := runtime.RunToolCommandResult(
		context.RepoRoot,
		context.ToolEnvironment,
		tool,
		capability,
		arguments...,
	)
	if runErr != nil && execution.FindingExitCode != 0 {
		var cmdErr runtime.CommandError
		if errors.As(runErr, &cmdErr) && cmdErr.Result.ExitCode == execution.FindingExitCode {
			return style.ExecutionResult{
				Output:  commandResult.Output,
				Command: runtime.BuildStyleCommandResult(commandResult),
			}, nil
		}
	}
	return style.ExecutionResult{
		Output:  commandResult.Output,
		Command: runtime.BuildStyleCommandResult(commandResult),
	}, runErr
}

func errUnknownTool(toolID string) (err error) {
	return errors.New("unknown tool " + toolID)
}
