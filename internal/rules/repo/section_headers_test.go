package repostyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

func TestCheckSectionHeadersFindsMissingHeaderInLongFile(t *testing.T) {
	repoRoot := t.TempDir()
	var builder strings.Builder
	builder.WriteString("package example\n\n")
	for range 105 {
		builder.WriteString("const value = 1\n")
	}

	fixtures.WriteFile(t, repoRoot, "internal/example/example.go", builder.String())

	output, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err == nil {
		t.Fatal("expected section header failure")
	}

	if !strings.Contains(output, "missing section headers") {
		t.Fatalf("expected missing-header message, got:\n%s", output)
	}
}

func TestCheckSectionHeadersAcceptsValidGoHeader(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 44) + " Types " + strings.Repeat("-", 43) + " */"
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\ntype Thing struct{}\n",
	)

	output, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err != nil {
		t.Fatalf("expected valid header to pass, output:\n%s", output)
	}
}

func TestCheckSectionHeadersCountsTabsAsFourColumns(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 43) + "\tTypes " + strings.Repeat("-", 43) + " */"
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\ntype Thing struct{}\n",
	)

	output, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err != nil {
		t.Fatalf("expected tab-width-aware header to pass, output:\n%s", output)
	}
}

/* ---------------------------------------- Header Names ---------------------------------------- */

func TestCheckSectionHeaderNamesFindsGenericHeadings(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 43) + " Checks " + strings.Repeat("-", 42) + " */"
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\nfunc run() {}\n",
	)

	output, err := CheckSectionHeaderNames(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err == nil {
		t.Fatal("expected generic section header failure")
	}

	if !strings.Contains(output, `generic section header name "Checks"`) {
		t.Fatalf("expected generic-header message, got:\n%s", output)
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

	output, err := CheckSectionHeaderNames(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err != nil {
		t.Fatalf("expected structural header to pass, output:\n%s", output)
	}
}
