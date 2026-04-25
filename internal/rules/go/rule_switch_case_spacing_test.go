package gostyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
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

	output, err := CheckSwitchCaseSpacing(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err == nil {
		t.Fatal("expected non-trivial switch spacing failure")
	}

	if !strings.Contains(
		output,
		"internal/example/example.go:7 non-trivial switch statements should separate "+
			"case blocks with a blank line",
	) {
		t.Fatalf("expected cramped switch violation, got:\n%s", output)
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

	output, err := CheckSwitchCaseSpacing(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err != nil {
		t.Fatalf("expected spaced non-trivial switch to pass, output:\n%s", output)
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

	output, err := CheckSwitchCaseSpacing(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err == nil {
		t.Fatal("expected very small switch spacing failure")
	}

	if !strings.Contains(
		output,
		"internal/example/example.go:8 very small switch statements should stay compact "+
			"without blank lines between case blocks",
	) {
		t.Fatalf("expected compact-switch violation, got:\n%s", output)
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

	output, err := CheckSwitchCaseSpacing(
		repoRoot,
		profiles.RepositoryConfig(t),
		contract.ScopeApp,
	)
	if err != nil {
		t.Fatalf("expected compact very small switch to pass, output:\n%s", output)
	}
}
