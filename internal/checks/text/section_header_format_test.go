package text

import (
	"strings"
	"testing"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
)

func TestCheckSectionHeadersFindsMissingHeaderInLongFile(t *testing.T) {
	repoRoot := t.TempDir()
	var builder strings.Builder
	builder.WriteString("package example\n\n")
	for range 105 {
		builder.WriteString("const value = 1\n")
	}

	testutil.WriteFile(t, repoRoot, "internal/example/example.go", builder.String())

	result, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		currentSectionHeaders(t),
		style.Scope("app"),
	)
	if err == nil {
		t.Fatal("expected section header failure")
	}

	if !hasDiagnostic(result, "text/section-headers/missing", "", 0, "missing section headers") {
		t.Fatalf("expected missing-header diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckSectionHeadersAcceptsValidGoHeader(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 44) + " Types " + strings.Repeat("-", 43) + " */"
	testutil.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\ntype Thing struct{}\n",
	)

	result, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		currentSectionHeaders(t),
		style.Scope("app"),
	)
	if err != nil {
		t.Fatalf("expected valid header to pass, diagnostics: %#v", result.Diagnostics)
	}
}

func TestCheckSectionHeadersCountsTabsAsFourColumns(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 43) + "\tTypes " + strings.Repeat("-", 43) + " */"
	testutil.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\ntype Thing struct{}\n",
	)

	result, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		currentSectionHeaders(t),
		style.Scope("app"),
	)
	if err != nil {
		t.Fatalf("expected tab-width-aware header to pass, diagnostics: %#v", result.Diagnostics)
	}
}
