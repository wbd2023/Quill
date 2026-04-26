package runner

import (
	"path/filepath"

	"ciphera/tools/internal/contract"
)

func FileCommandArguments(
	repoRoot string,
	spec contract.ExecutionSpec,
) (arguments []string) {
	detail, found := spec.FileCommandExecution()
	if !found {
		return nil
	}

	arguments = append([]string{}, detail.Arguments...)
	if detail.ConfigFile != "" {
		arguments = append(
			arguments,
			detail.ConfigArgument,
			filepath.Join(repoRoot, detail.ConfigFile),
		)
	}

	return arguments
}
