package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

/* -------------------------------------- Import Boundaries ------------------------------------- */

func TestStylePlatformImportBoundaries(t *testing.T) {
	t.Parallel()

	toolsRoot := importBoundaryRoot(t)
	testCases := []struct {
		name      string
		directory string
		recursive bool
		forbidden []string
	}{
		{
			name:      "contract does not import internal packages",
			directory: "internal/contract",
			forbidden: []string{"ciphera/tools/internal/"},
		},
		{
			name:      "profile depends only on contracts and policy",
			directory: "internal/profile",
			forbidden: []string{
				"ciphera/tools/internal/executors",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules",
			},
		},
		{
			name:      "policy depends only on contracts",
			directory: "internal/policy",
			forbidden: []string{
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/executors",
				"ciphera/tools/internal/filewalk",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/styleguide",
			},
		},
		{
			name:      "runner is generic execution machinery",
			directory: "internal/runner",
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runtime",
			},
		},
		{
			name:      "toolchain depends only on contracts",
			directory: "internal/toolchain",
			forbidden: []string{
				"ciphera/tools/internal/architecture",
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/executors",
				"ciphera/tools/internal/filewalk",
				"ciphera/tools/internal/policy",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/styleguide",
			},
		},
		{
			name:      "rulepack depends only on contracts and toolchain",
			directory: "internal/rulepack",
			forbidden: []string{
				"ciphera/tools/internal/architecture",
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/executors",
				"ciphera/tools/internal/filewalk",
				"ciphera/tools/internal/policy",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/styleguide",
			},
		},
		{
			name:      "filewalk does not import profile",
			directory: "internal/filewalk",
			forbidden: []string{"ciphera/tools/internal/profile"},
		},
		{
			name:      "styleguide does not import profile or rulepack",
			directory: "internal/styleguide",
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/rulepack",
			},
		},
		{
			name:      "go checks do not import rulepack",
			directory: "internal/rules/golang/checks",
			forbidden: []string{"ciphera/tools/internal/rulepack"},
		},
		{
			name:      "go order checks do not import rulepack",
			directory: "internal/rules/golang/order",
			forbidden: []string{"ciphera/tools/internal/rulepack"},
		},
		{
			name:      "bash checks use filewalk directly",
			directory: "internal/rules/bash",
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/rules/text",
			},
		},
		{
			name:      "go checks do not depend on text helpers",
			directory: "internal/rules/golang",
			recursive: true,
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/rules/text",
			},
		},
		{
			name:      "text checks do not import profile",
			directory: "internal/rules/text",
			forbidden: []string{"ciphera/tools/internal/profile"},
		},
		{
			name:      "security checks do not import profile",
			directory: "internal/rules/security",
			forbidden: []string{"ciphera/tools/internal/profile"},
		},
		{
			name:      "naming checks do not import profile",
			directory: "internal/rules/naming",
			forbidden: []string{"ciphera/tools/internal/profile"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			directory := filepath.Join(toolsRoot, testCase.directory)
			files := productionGoFiles(t, directory, testCase.recursive)
			for _, file := range files {
				for _, imported := range fileImports(t, file) {
					if !forbiddenImport(imported, testCase.forbidden) {
						continue
					}

					t.Fatalf("%s imports forbidden package %q", file, imported)
				}
			}
		})
	}
}

func TestRetiredHelperFilesStayRetired(t *testing.T) {
	t.Parallel()

	toolsRoot := importBoundaryRoot(t)
	retired := []string{
		"internal/rules/bash/results.go",
		"internal/rules/text/results.go",
		"internal/rules/security/results.go",
		"internal/rules/naming/results.go",
		"internal/rules/golang/scenarios/behaviour_harness_test.go",
	}

	for _, path := range retired {
		if _, err := os.Stat(filepath.Join(toolsRoot, path)); err == nil {
			t.Fatalf("retired helper file still exists: %s", path)
		}
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func importBoundaryRoot(t *testing.T) (toolsRoot string) {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve import-boundary test path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
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

func forbiddenImport(imported string, forbidden []string) (found bool) {
	for _, forbiddenPath := range forbidden {
		if imported == forbiddenPath || strings.HasPrefix(imported, forbiddenPath) {
			return true
		}
	}

	return false
}
