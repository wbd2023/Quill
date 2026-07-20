package text

import (
	"strings"
	"testing"

	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func TestCheckSectionHeaderDensityWarnsForShortFileHeader(t *testing.T) {
	repoRoot := t.TempDir()
	header := "/* " + strings.Repeat("-", 42) + " Helpers " + strings.Repeat("-", 42) + " */"
	testutil.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		"package example\n\n"+header+"\n\nfunc run() {}\n",
	)

	result, _ := CheckSectionHeaderDensity(
		repoRoot,
		profiles.RepositoryConfig(t),
		currentSectionHeaders(t),
		style.Scope("all"),
	)
	if len(result.Diagnostics) == 0 {
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
	testutil.WriteFile(t, repoRoot, "internal/example/example.go", builder.String())

	result, err := CheckSectionHeaderDensity(
		repoRoot,
		profiles.RepositoryConfig(t),
		currentSectionHeaders(t),
		style.Scope("all"),
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
	for range 7 {
		builder.WriteString(header + "\n\n")
	}
	testutil.WriteFile(t, repoRoot, "internal/example/example.go", builder.String())

	result, _ := CheckSectionHeaderDensity(
		repoRoot,
		profiles.RepositoryConfig(t),
		currentSectionHeaders(t),
		style.Scope("all"),
	)
	if len(result.Diagnostics) == 0 {
		t.Fatal("expected over-dense header warning")
	}

	if !hasDiagnostic(
		result,
		"text/section-header-density/too-many",
		"internal/example/example.go",
		0,
		"7 section headers",
	) {
		t.Fatalf("expected too-many diagnostic, got: %#v", result.Diagnostics)
	}
}
