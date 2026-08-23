package architecture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestProductionSourceStaysWithinEntrypointAndInternal enforces ADR 0004: Quill is CLI-first
// and language neutral, so it ships no supported in-process Go library. Production source may
// exist only in the cmd/quill entrypoint or a private internal package.
func TestProductionSourceStaysWithinEntrypointAndInternal(t *testing.T) {
	t.Parallel()

	repositoryRoot := importBoundaryRoot(t)
	for _, file := range productionGoFiles(t, repositoryRoot, true, nil) {
		relative, err := filepath.Rel(repositoryRoot, file)
		if err != nil {
			t.Fatalf("rel %s: %v", file, err)
		}

		relative = filepath.ToSlash(relative)
		if strings.HasPrefix(relative, "cmd/quill/") ||
			strings.HasPrefix(relative, "internal/") {
			continue
		}

		t.Fatalf(
			"production Go file %s lives outside cmd/quill or internal; "+
				"all implementation packages must be private per ADR 0004",
			relative,
		)
	}
}

// TestNoSourceImportsRootPackage enforces that no Go file imports the bare module path, which
// was the removed public facade. The architecture layer classifier deliberately ignores the
// bare module path, so this guard keeps that contract meaningful.
func TestNoSourceImportsRootPackage(t *testing.T) {
	t.Parallel()

	toolsRoot := importBoundaryRoot(t)
	forbiddenRoot := moduleImportPath(t, toolsRoot)

	err := filepath.WalkDir(
		toolsRoot,
		func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if entry.IsDir() {
				switch entry.Name() {
				case "testdata", ".cache", ".git", "vendor", "third_party":
					return filepath.SkipDir
				}
				return nil
			}

			if filepath.Ext(path) != ".go" {
				return nil
			}

			for _, imported := range fileImports(t, path) {
				if imported == forbiddenRoot {
					relative, relErr := filepath.Rel(toolsRoot, path)
					if relErr != nil {
						t.Fatalf("rel %s: %v", path, relErr)
					}

					t.Fatalf(
						"%s imports the removed root package %q; "+
							"route through internal/engine per ADR 0004",
						filepath.ToSlash(relative),
						forbiddenRoot,
					)
				}
			}

			return nil
		},
	)
	if err != nil {
		t.Fatalf("walk repository: %v", err)
	}
}
