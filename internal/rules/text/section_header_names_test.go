package text

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

func TestCheckSectionHeaderNamesFindsGenericHeadings(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 43) + " Checks " + strings.Repeat("-", 42) + " */"
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\nfunc run() {}\n",
	)

	result, err := CheckSectionHeaderNames(
		repoRoot,
		profiles.RepositoryConfig(t),
		currentSectionHeaders(t),
		contract.Scope("app"),
	)
	if err == nil {
		t.Fatal("expected generic section header failure")
	}

	if !hasDiagnostic(
		result,
		"text/section-header-names/generic",
		"internal/example/example.go",
		3,
		`generic section header name "Checks"`,
	) {
		t.Fatalf("expected generic-header diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckSectionHeaderNamesAllowsStructuralHeadings(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 42) + " Helpers " + strings.Repeat("-", 42) + " */"
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\nfunc run() {}\n",
	)

	result, err := CheckSectionHeaderNames(
		repoRoot,
		profiles.RepositoryConfig(t),
		currentSectionHeaders(t),
		contract.Scope("app"),
	)
	if err != nil {
		t.Fatalf("expected structural header to pass, diagnostics: %#v", result.Diagnostics)
	}
}

func TestCheckSectionHeaderNamesUsesProfileGenericNames(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 43) + " Local " + strings.Repeat("-", 43) + " */"
	headers := currentSectionHeaders(t)
	headers.GenericNames = []string{"Local"}
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\nfunc run() {}\n",
	)

	result, err := CheckSectionHeaderNames(
		repoRoot,
		profiles.RepositoryConfig(t),
		headers,
		contract.Scope("app"),
	)
	if err == nil {
		t.Fatal("expected profile-owned generic section header failure")
	}

	if !hasDiagnostic(
		result,
		"text/section-header-names/generic",
		"internal/example/example.go",
		3,
		`generic section header name "Local"`,
	) {
		t.Fatalf("expected profile generic-header diagnostic, got: %#v", result.Diagnostics)
	}
}
