package execution

import (
	"path/filepath"

	"github.com/wbd2023/Quill/internal/style"
)

// FileCommandArguments extracts the command arguments from a file-command job, resolves its config
// file path against the repository root, and appends the selected files.
func FileCommandArguments(
	repoRoot string,
	job style.Job,
	files []string,
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

	arguments = append(arguments, files...)

	return arguments
}
