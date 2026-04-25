package runner

import (
	"errors"
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runtime"
)

func FileCommandExecutor(
	context Context,
	spec contract.ExecutionSpec,
	_ map[string]runtime.ToolStatus,
) (output string, err error) {
	files, err := CollectFileSetFiles(context, spec.FileSet)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", nil
	}

	tool, found := context.Effective.ToolByID(spec.ToolID)
	if !found {
		return "", errUnknownTool(spec.ToolID)
	}

	arguments := FileCommandArguments(context.RepoRoot, spec)
	arguments = append(arguments, files...)
	return runtime.RunCommand(context.RepoRoot, context.ToolEnvironment, tool.Command, arguments...)
}

func FileCommandArguments(
	repoRoot string,
	spec contract.ExecutionSpec,
) (arguments []string) {
	arguments = append([]string{}, spec.Arguments...)
	if spec.ConfigFile != "" {
		arguments = append(arguments, spec.ConfigArgument, filepath.Join(repoRoot, spec.ConfigFile))
	}

	return arguments
}

func errUnknownTool(toolID string) (err error) {
	return errors.New("unknown tool " + toolID)
}
