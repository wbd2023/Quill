package architecture

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNoPublicProductionPackageAtRoot enforces ADR 0004: Quill is CLI-first and language
// neutral, so it ships no supported in-process Go library. No production Go file may live in
// the repository root, and no source file may import the bare module path (the removed root
// facade). All production implementation packages must live under internal/.
func TestNoPublicProductionPackageAtRoot(t *testing.T) {
	t.Parallel()

	toolsRoot := importBoundaryRoot(t)

	entries, err := os.ReadDir(toolsRoot)
	if err != nil {
		t.Fatalf("read repository root: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(toolsRoot, entry.Name())
		if !isProductionGoFile(path) {
			continue
		}

		relative, err := filepath.Rel(toolsRoot, path)
		if err != nil {
			t.Fatalf("rel %s: %v", path, err)
		}

		t.Fatalf(
			"production Go file %s lives at the repository root; "+
				"all production packages must live under internal/ per ADR 0004",
			filepath.ToSlash(relative),
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
				case "testdata", ".git", "vendor", "third_party":
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
