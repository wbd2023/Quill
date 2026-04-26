package text

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

	result, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		profiles.Current(t).Formatting.SectionHeaders,
		contract.Scope("app"),
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
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\ntype Thing struct{}\n",
	)

	result, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		profiles.Current(t).Formatting.SectionHeaders,
		contract.Scope("app"),
	)
	if err != nil {
		t.Fatalf("expected valid header to pass, diagnostics: %#v", result.Diagnostics)
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

	result, err := CheckSectionHeaders(
		repoRoot,
		profiles.RepositoryConfig(t),
		profiles.Current(t).Formatting.SectionHeaders,
		contract.Scope("app"),
	)
	if err != nil {
		t.Fatalf("expected tab-width-aware header to pass, diagnostics: %#v", result.Diagnostics)
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

	result, err := CheckSectionHeaderNames(
		repoRoot,
		profiles.RepositoryConfig(t),
		profiles.Current(t).Formatting.SectionHeaders,
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
		profiles.Current(t).Formatting.SectionHeaders,
		contract.Scope("app"),
	)
	if err != nil {
		t.Fatalf("expected structural header to pass, diagnostics: %#v", result.Diagnostics)
	}
}

func TestCheckSectionHeaderNamesUsesProfileGenericNames(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 43) + " Local " + strings.Repeat("-", 43) + " */"
	config := profiles.Current(t)
	config.Formatting.SectionHeaders.GenericNames = []string{"Local"}
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\nfunc run() {}\n",
	)

	result, err := CheckSectionHeaderNames(
		repoRoot,
		profiles.RepositoryConfig(t),
		config.Formatting.SectionHeaders,
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

func TestCheckSectionHeaderDensityWarnsForShortFileHeader(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 42) + " Helpers " + strings.Repeat("-", 42) + " */"
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\nfunc run() {}\n",
	)

	result, err := CheckSectionHeaderDensity(
		repoRoot,
		profiles.RepositoryConfig(t),
		profiles.Current(t).Formatting.SectionHeaders,
		contract.Scope("app"),
	)
	if err == nil {
		t.Fatal("expected short-file header density warning")
	}

	if !hasDiagnostic(
		result,
		"text/section-header-density/short-file",
		"internal/example/example.go",
		0,
		"short 5-line file has section headers",
	) {
		t.Fatalf("expected short-file diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckSectionHeaderDensityAllowsEightyLineFile(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 42) + " Helpers " + strings.Repeat("-", 42) + " */"
	var builder strings.Builder
	builder.WriteString("package example\n\n")
	builder.WriteString(header + "\n")
	for range 77 {
		builder.WriteString("const value = 1\n")
	}
	fixtures.WriteFile(t, repoRoot, "internal/example/example.go", builder.String())

	result, err := CheckSectionHeaderDensity(
		repoRoot,
		profiles.RepositoryConfig(t),
		profiles.Current(t).Formatting.SectionHeaders,
		contract.Scope("app"),
	)
	if err != nil {
		t.Fatalf(
			"expected 80-line file to pass density check, diagnostics: %#v",
			result.Diagnostics,
		)
	}
}

func TestCheckSectionHeaderDensityWarnsForManyHeaders(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 42) + " Helpers " + strings.Repeat("-", 42) + " */"
	var builder strings.Builder
	builder.WriteString("package example\n\n")
	for range 6 {
		builder.WriteString(header + "\n\n")
	}
	fixtures.WriteFile(t, repoRoot, "internal/example/example.go", builder.String())

	result, err := CheckSectionHeaderDensity(
		repoRoot,
		profiles.RepositoryConfig(t),
		profiles.Current(t).Formatting.SectionHeaders,
		contract.Scope("app"),
	)
	if err == nil {
		t.Fatal("expected over-dense header warning")
	}

	if !hasDiagnostic(
		result,
		"text/section-header-density/too-many",
		"internal/example/example.go",
		0,
		"6 section headers",
	) {
		t.Fatalf("expected too-many diagnostic, got: %#v", result.Diagnostics)
	}
}
