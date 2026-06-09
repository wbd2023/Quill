package runner

import (
	"path/filepath"

	"ciphera/tools/internal/style"
)

func FileCommandArguments(
	repoRoot string,
	spec style.ExecutionSpec,
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
