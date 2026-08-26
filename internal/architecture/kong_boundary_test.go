package architecture

import (
	"path/filepath"
	"testing"
)

func TestKongImportBoundary(t *testing.T) {
	t.Parallel()

	root := importBoundaryRoot(t)
	cliDirectory := filepath.Join(root, "internal", "cli")
	sourceRoots := []string{
		filepath.Join(root, "cmd"),
		filepath.Join(root, "internal"),
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
