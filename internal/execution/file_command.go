package execution

import (
	"path/filepath"

	"github.com/wbd2023/quill/internal/style"
)

// FileCommandArguments extracts the command arguments from a file-command Job, resolves its
// config file path against the repository root, and appends the selected files.
func FileCommandArguments(
	repoRoot string,
	command style.FileCommand,
	files []string,
) (arguments []string) {
	arguments = append([]string{}, command.Arguments...)
	if command.ConfigFile != "" {
		arguments = append(
			arguments,
			command.ConfigArgument,
			filepath.Join(repoRoot, command.ConfigFile),
		)
	}

	arguments = append(arguments, files...)

	return arguments
}
