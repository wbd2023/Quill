package architecture

import (
	"path/filepath"
	"testing"
)

func TestKongImportBoundary(t *testing.T) {
	t.Parallel()

	repositoryRoot := importBoundaryRoot(t)
	cliDirectory := filepath.Join(repositoryRoot, "internal", "cli")
	sourceRoots := []string{
		filepath.Join(repositoryRoot, "cmd"),
		filepath.Join(repositoryRoot, "internal"),
	}

	for _, sourceRoot := range sourceRoots {
		files := productionGoFiles(t, sourceRoot, true, nil)
		for _, file := range files {
			for _, imported := range fileImports(t, file) {
				if imported != "github.com/alecthomas/kong" {
					continue
				}
				if filepath.Dir(file) == cliDirectory {
					continue
				}

				t.Fatalf("%s imports Kong outside internal/cli", file)
			}
		}
	}
}
