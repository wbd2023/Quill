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
	for _, testCase := range importBoundaryCases() {
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

func TestRetiredPathsStayRetired(t *testing.T) {
	t.Parallel()

	toolsRoot := importBoundaryRoot(t)
	retired := []string{
		"internal/rules",
		"internal/rules/bash/results.go",
		"internal/rules/text/results.go",
		"internal/rules/security/results.go",
		"internal/rules/vocabulary/results.go",
		"internal/rules/golang/check_ids.go",
		"internal/rules/golang/checks",
		"internal/rules/golang/diagnostic_test.go",
		"internal/rules/golang/order",
		"internal/rules/golang/rule_architecture.go",
		"internal/rules/golang/rule_guard_clause_spacing.go",
		"internal/rules/golang/rule_guard_clause_spacing_test.go",
		"internal/rules/golang/rule_switch_case_spacing.go",
		"internal/rules/golang/rule_switch_case_spacing_test.go",
		"internal/rules/golang/scenarios/behaviour_harness_test.go",
		"internal/rules/golang/scenarios/domain_identifier_casts_test.go",
		"internal/pack/shipped/go_target_ids.go",
		"internal/pack/shipped/pack_ids.go",
		"internal/pack/shipped/project_check_ids.go",
		"internal/pack/shipped/scanner_ids.go",
		"internal/runtime/handlers_test.go",
		"internal/runtime/tool_inspection.go",
		"internal/runtime/tool_version.go",
		"internal/runtime/tool_version_detection.go",
		"internal/runtime/version_normalisation.go",
	}

	for _, path := range retired {
		if _, err := os.Stat(filepath.Join(toolsRoot, path)); err == nil {
			t.Fatalf("retired path still exists: %s", path)
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
