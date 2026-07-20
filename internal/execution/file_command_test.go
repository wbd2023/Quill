package execution

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/wbd2023/Quill/internal/style"
)

func TestFileCommandArgumentsAppendsSelectedFiles(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	files := []string{
		filepath.Join(repoRoot, "README.md"),
		filepath.Join(repoRoot, "docs", "architecture.md"),
	}
	job := style.FileCommandExecution{
		Arguments:      []string{"--check"},
		ConfigArgument: "--config",
		ConfigFile:     ".markdownlint.jsonc",
	}

	arguments := FileCommandArguments(repoRoot, job, files)
	want := []string{
		"--check",
		"--config",
		filepath.Join(repoRoot, ".markdownlint.jsonc"),
		files[0],
		files[1],
	}
	if !slices.Equal(arguments, want) {
		t.Fatalf("FileCommandArguments() = %q, want %q", arguments, want)
	}
}
