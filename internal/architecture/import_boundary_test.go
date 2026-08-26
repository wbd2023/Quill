package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/testutil"
)

/* -------------------------------------- Import Boundaries ------------------------------------- */

func TestStylePlatformImportBoundaries(t *testing.T) {
	t.Parallel()

	root := importBoundaryRoot(t)
	modulePath := moduleImportPath(t, root)
	for _, testCase := range importBoundaryCases() {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			directory := filepath.Join(root, testCase.directory)
			files := productionGoFiles(t, directory, testCase.recursive, testCase.excludeSubdirs)
			for _, file := range files {
				for _, imported := range fileImports(t, file) {
					if !hasForbiddenImportExcept(
						imported, modulePath, testCase.allowed, testCase.forbidden,
					) {
						continue
					}

					t.Fatalf("%s imports forbidden package %q", file, imported)
				}
			}
		})
	}
}

func TestHasForbiddenImportExceptMatchesWholePackagePaths(t *testing.T) {
	t.Parallel()

	const modulePath = "example.test/quill"
	testCases := []struct {
		name      string
		imported  string
		forbidden []string
		want      bool
	}{
		{
			name:      "exact package",
			imported:  modulePath + "/internal/process",
			forbidden: []string{"internal/process"},
			want:      true,
		},
		{
			name:      "package child",
			imported:  modulePath + "/internal/process/runner",
			forbidden: []string{"internal/process"},
			want:      true,
		},
		{
			name:      "similarly prefixed sibling",
			imported:  modulePath + "/internal/processors",
			forbidden: []string{"internal/process"},
			want:      false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := hasForbiddenImportExcept(
				testCase.imported,
				modulePath,
				nil,
				testCase.forbidden,
			)
			if got != testCase.want {
				t.Fatalf("hasForbiddenImportExcept() = %t, want %t", got, testCase.want)
			}
		})
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func importBoundaryRoot(t *testing.T) (root string) {
	t.Helper()
	return testutil.RepositoryRoot(t)
}

func moduleImportPath(t *testing.T, root string) (modulePath string) {
	t.Helper()

	contents, err := os.ReadFile(filepath.Join(root, "go.mod"))
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
	excludeSubdirs []string,
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
					if path == directory {
						return nil
					}
					switch entry.Name() {
					case "testdata", ".cache", ".git", "vendor", "third_party":
						return filepath.SkipDir
					}
					if isExcludedSubdir(path, directory, excludeSubdirs) {
						return filepath.SkipDir
					}
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

// isExcludedSubdir reports whether path is a direct child of root whose basename appears in
// excludeSubdirs. It is used to keep a recursive boundary check strict while admitting a named
// child package that is allowed to cross the parent's import line.
func isExcludedSubdir(path string, root string, excludeSubdirs []string) (excluded bool) {
	if len(excludeSubdirs) == 0 {
		return false
	}

	parent := filepath.Dir(path)
	if filepath.Clean(parent) != filepath.Clean(root) {
		return false
	}

	name := filepath.Base(path)
	for _, excluded := range excludeSubdirs {
		if name == excluded {
			return true
		}
	}

	return false
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

func hasForbiddenImport(imported string, modulePath string, forbidden []string) (found bool) {
	return hasForbiddenImportExcept(imported, modulePath, nil, forbidden)
}

func hasForbiddenImportExcept(
	imported string,
	modulePath string,
	allowed []string,
	forbidden []string,
) (found bool) {
	localPrefix := modulePath + "/"
	if !strings.HasPrefix(imported, localPrefix) {
		return false
	}

	relative := strings.TrimPrefix(imported, localPrefix)
	for _, allowedPath := range allowed {
		if relative == allowedPath || strings.HasPrefix(relative, allowedPath+"/") {
			return false
		}
	}

	for _, forbiddenPath := range forbidden {
		if relative == forbiddenPath || strings.HasPrefix(relative, forbiddenPath+"/") {
			return true
		}
	}

	return false
}
