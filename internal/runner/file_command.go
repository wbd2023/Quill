package runner

import (
	"path/filepath"

	"ciphera/tools/internal/contract"
)

func FileCommandArguments(
	repoRoot string,
	spec contract.ExecutionSpec,
) (arguments []string) {
	execution, found := spec.FileCommandExecution()
	if !found {
		return nil
	}

	arguments = append([]string{}, execution.Arguments...)
	if execution.ConfigFile != "" {
		arguments = append(
			arguments,
			execution.ConfigArgument,
			filepath.Join(repoRoot, execution.ConfigFile),
		)
	}

	return arguments
}
