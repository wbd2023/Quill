package gostyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
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

	output, err := CheckGuardClauseSpacing(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err == nil {
		t.Fatal("expected guard-clause spacing failure")
	}

	if !strings.Contains(output, "internal/example/example.go:7") {
		t.Fatalf("expected output to include offending line, got:\n%s", output)
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

	output, err := CheckGuardClauseSpacing(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err != nil {
		t.Fatalf("expected spaced guard clauses to pass, output:\n%s", output)
	}
}
