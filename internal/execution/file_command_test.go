package execution

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/wbd2023/quill/internal/style"
)

func TestFileCommandArgumentsAppendsSelectedFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	files := []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "docs", "architecture.md"),
	}
	job := style.FileCommand{
		Arguments:      []string{"--check"},
		ConfigArgument: "--config",
		ConfigFile:     ".markdownlint.jsonc",
	}

	arguments := FileCommandArguments(root, job, files)
	want := []string{
		"--check",
		"--config",
		filepath.Join(root, ".markdownlint.jsonc"),
		files[0],
		files[1],
	}
	if !slices.Equal(arguments, want) {
		t.Fatalf("FileCommandArguments() = %q, want %q", arguments, want)
	}
}
