package runner

import (
	"path/filepath"

	"ciphera/tools/internal/style"
)

// FileCommandArguments extracts the command arguments from a file-command job, resolving config
// file paths against the repository root.
func FileCommandArguments(
	repoRoot string,
	job style.Job,
) (arguments []string) {
	execution, found := job.(style.FileCommandExecution)
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
