package golang

import (
	"path/filepath"
	"strings"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/rules/golang/checks"
)

func TestCheckSwitchCaseSpacingFindsCrampedNonTrivialSwitches(t *testing.T) {
	repoRoot := t.TempDir()
	source := strings.Join([]string{
		"package example",
		"",
		"func Render(value string) string {",
		"\tswitch value {",
		"\tcase \"a\":",
		"\t\treturn \"A\"",
		"\tcase \"b\":",
		"\t\treturn \"B\"",
		"\tcase \"c\":",
		"\t\treturn \"C\"",
		"\tdefault:",
		"\t\treturn \"?\"",
		"\t}",
		"}",
		"",
	}, "\n")
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		source,
	)

	result, err := CheckSwitchCaseSpacing(
		repoRoot,
		[]string{filepath.Join(repoRoot, "internal")},
		profiles.RepositoryConfig(t),
	)
	if err == nil {
		t.Fatal("expected non-trivial switch spacing failure")
	}

	if !hasDiagnostic(
		result,
		checks.DiagnosticSwitchCaseSpacing,
		"internal/example/example.go",
		7,
		"non-trivial switch statements should separate case blocks",
	) {
		t.Fatalf("expected cramped switch diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckSwitchCaseSpacingAllowsSeparatedNonTrivialSwitches(t *testing.T) {
	repoRoot := t.TempDir()
	source := strings.Join([]string{
		"package example",
		"",
		"func Render(value string) string {",
		"\tswitch value {",
		"\tcase \"a\":",
		"\t\treturn \"A\"",
		"",
		"\tcase \"b\":",
		"\t\treturn \"B\"",
		"",
		"\tcase \"c\":",
		"\t\treturn \"C\"",
		"",
		"\tdefault:",
		"\t\treturn \"?\"",
		"\t}",
		"}",
		"",
	}, "\n")
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		source,
	)

	result, err := CheckSwitchCaseSpacing(
		repoRoot,
		[]string{filepath.Join(repoRoot, "internal")},
		profiles.RepositoryConfig(t),
	)
	if err != nil {
		t.Fatalf("expected spaced non-trivial switch to pass, diagnostics: %#v", result.Diagnostics)
	}
}

/* -------------------------------------- Compact Switches -------------------------------------- */

func TestCheckSwitchCaseSpacingRejectsOverSpacedVerySmallSwitches(t *testing.T) {
	repoRoot := t.TempDir()
	source := strings.Join([]string{
		"package example",
		"",
		"func Render(value string) string {",
		"\tswitch value {",
		"\tcase \"a\":",
		"\t\treturn \"A\"",
		"",
		"\tdefault:",
		"\t\treturn \"?\"",
		"\t}",
		"}",
		"",
	}, "\n")
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		source,
	)

	result, err := CheckSwitchCaseSpacing(
		repoRoot,
		[]string{filepath.Join(repoRoot, "internal")},
		profiles.RepositoryConfig(t),
	)
	if err == nil {
		t.Fatal("expected very small switch spacing failure")
	}

	if !hasDiagnostic(
		result,
		checks.DiagnosticSwitchCaseSpacing,
		"internal/example/example.go",
		8,
		"very small switch statements should stay compact",
	) {
		t.Fatalf("expected compact-switch diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckSwitchCaseSpacingAllowsCompactVerySmallSwitches(t *testing.T) {
	repoRoot := t.TempDir()
	source := strings.Join([]string{
		"package example",
		"",
		"func Render(value string) string {",
		"\tswitch value {",
		"\tcase \"a\":",
		"\t\treturn \"A\"",
		"\tdefault:",
		"\t\treturn \"?\"",
		"\t}",
		"}",
		"",
	}, "\n")
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/example.go",
		source,
	)

	result, err := CheckSwitchCaseSpacing(
		repoRoot,
		[]string{filepath.Join(repoRoot, "internal")},
		profiles.RepositoryConfig(t),
	)
	if err != nil {
		t.Fatalf("expected compact very small switch to pass, diagnostics: %#v", result.Diagnostics)
	}
}
