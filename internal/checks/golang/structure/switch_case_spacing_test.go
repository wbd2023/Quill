package structure

import (
	"strings"
	"testing"

	"ciphera/tools/internal/checks/golang/analysis"
)

func TestCheckSwitchCaseSpacingFindsCrampedNonTrivialSwitches(t *testing.T) {
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
	fileSet, file := parseGoSource(t, source)

	violations := CheckSwitchCaseSpacing(fileSet, file, sourceLines(source))
	if !hasViolationAt(
		violations,
		analysis.DiagnosticSwitchCaseSpacing,
		7,
		"non-trivial switch statements should separate case blocks",
	) {
		t.Fatalf("expected cramped switch violation, got: %#v", violations)
	}
}

func TestCheckSwitchCaseSpacingAllowsSeparatedNonTrivialSwitches(t *testing.T) {
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
	fileSet, file := parseGoSource(t, source)

	violations := CheckSwitchCaseSpacing(fileSet, file, sourceLines(source))
	if len(violations) != 0 {
		t.Fatalf("expected no cramped switch violations, got: %#v", violations)
	}
}

/* --------------------------------------- Small Switches --------------------------------------- */

func TestCheckSwitchCaseSpacingAllowsSpacedSmallSwitches(t *testing.T) {
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
	fileSet, file := parseGoSource(t, source)

	violations := CheckSwitchCaseSpacing(fileSet, file, sourceLines(source))
	if len(violations) != 0 {
		t.Fatalf("expected no violations for spaced small switch, got: %#v", violations)
	}
}

func TestCheckSwitchCaseSpacingAllowsCompactSmallSwitches(t *testing.T) {
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
	fileSet, file := parseGoSource(t, source)

	violations := CheckSwitchCaseSpacing(fileSet, file, sourceLines(source))
	if len(violations) != 0 {
		t.Fatalf("expected no violations for compact small switch, got: %#v", violations)
	}
}
