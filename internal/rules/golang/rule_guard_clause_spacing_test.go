package golang

import (
	"path/filepath"
	"strings"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/rules/golang/checks"
)

func TestCheckGuardClauseSpacingFindsAdjacentGuardClauses(t *testing.T) {
	repoRoot := t.TempDir()
	source := strings.Join([]string{
		"package example",
		"",
		"func Validate(a int, b int) error {",
		"\tif a == 0 {",
		"\t\treturn nil",
		"\t}",
		"\tif b == 0 {",
		"\t\treturn nil",
		"\t}",
		"",
		"\treturn nil",
		"}",
		"",
	}, "\n")
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		source,
	)

	result, err := CheckGuardClauseSpacing(
		repoRoot,
		[]string{filepath.Join(repoRoot, "internal")},
		profiles.RepositoryConfig(t),
	)
	if err == nil {
		t.Fatal("expected guard-clause spacing failure")
	}

	if !hasDiagnostic(
		result,
		checks.DiagnosticGuardClauseSpacing,
		"internal/example/example.go",
		7,
		"separated by a blank line",
	) {
		t.Fatalf("expected guard-clause diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckGuardClauseSpacingAllowsSeparatedGuardClauses(t *testing.T) {
	repoRoot := t.TempDir()
	source := strings.Join([]string{
		"package example",
		"",
		"func Validate(a int, b int) error {",
		"\tif a == 0 {",
		"\t\treturn nil",
		"\t}",
		"",
		"\tif b == 0 {",
		"\t\treturn nil",
		"\t}",
		"",
		"\treturn nil",
		"}",
		"",
	}, "\n")
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		source,
	)

	result, err := CheckGuardClauseSpacing(
		repoRoot,
		[]string{filepath.Join(repoRoot, "internal")},
		profiles.RepositoryConfig(t),
	)
	if err != nil {
		t.Fatalf("expected spaced guard clauses to pass, diagnostics: %#v", result.Diagnostics)
	}
}
