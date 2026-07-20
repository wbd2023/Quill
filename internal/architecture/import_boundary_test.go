package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wbd2023/Quill/internal/testutil"
)

/* -------------------------------------- Import Boundaries ------------------------------------- */

func TestStylePlatformImportBoundaries(t *testing.T) {
	t.Parallel()

	toolsRoot := importBoundaryRoot(t)
	modulePath := moduleImportPath(t, toolsRoot)
	for _, testCase := range importBoundaryCases() {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			directory := filepath.Join(toolsRoot, testCase.directory)
			files := productionGoFiles(t, directory, testCase.recursive)
			for _, file := range files {
				for _, imported := range fileImports(t, file) {
					if !forbiddenImport(imported, modulePath, testCase.forbidden) {
						continue
					}

					t.Fatalf("%s imports forbidden package %q", file, imported)
				}
			}
		})
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func importBoundaryRoot(t *testing.T) (toolsRoot string) {
	t.Helper()
	return testutil.RepositoryRoot(t)
}

func moduleImportPath(t *testing.T, repositoryRoot string) (modulePath string) {
	t.Helper()

	contents, err := os.ReadFile(filepath.Join(repositoryRoot, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}

	for _, line := range strings.Split(string(contents), "\n") {
		if modulePath, found := strings.CutPrefix(strings.TrimSpace(line), "module "); found {
			return modulePath
		}
	}

	t.Fatal("go.mod has no module directive")
	return ""
}

func productionGoFiles(
	t *testing.T,
	directory string,
	recursive bool,
) (files []string) {
	t.Helper()

	if recursive {
		err := filepath.WalkDir(
			directory,
			func(path string, entry os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}

				if entry.IsDir() {
					return nil
				}

				if isProductionGoFile(path) {
					files = append(files, path)
				}

				return nil
			},
		)
		if err != nil {
			t.Fatalf("walk %s: %v", directory, err)
		}

		return files
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read %s: %v", directory, err)
	}

	for _, entry := range entries {
		path := filepath.Join(directory, entry.Name())
		if entry.IsDir() || !isProductionGoFile(path) {
			continue
		}

		files = append(files, path)
	}

	return files
}

func isProductionGoFile(path string) (production bool) {
	return strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go")
}

func fileImports(t *testing.T, path string) (imports []string) {
	t.Helper()

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	for _, imported := range file.Imports {
		imports = append(imports, strings.Trim(imported.Path.Value, `"`))
	}

	return imports
}

func forbiddenImport(imported string, modulePath string, forbidden []string) (found bool) {
	localPrefix := modulePath + "/"
	if !strings.HasPrefix(imported, localPrefix) {
		return false
	}

	relative := strings.TrimPrefix(imported, localPrefix)
	for _, forbiddenPath := range forbidden {
		if relative == forbiddenPath || strings.HasPrefix(relative, forbiddenPath) {
			return true
		}
	}
	return false
}
