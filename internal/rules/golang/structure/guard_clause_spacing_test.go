package structure

import (
	"strings"
	"testing"

	"ciphera/tools/internal/rules/golang/analysis"
)

func TestCheckGuardClauseSpacingFindsAdjacentGuardClauses(t *testing.T) {
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
	fileSet, file := parseGoSource(t, source)

	violations := CheckGuardClauseSpacing(fileSet, file)
	if !hasViolationAt(
		violations,
		analysis.DiagnosticGuardClauseSpacing,
		7,
		"separated by a blank line",
	) {
		t.Fatalf("expected guard-clause violation, got: %#v", violations)
	}
}

func TestCheckGuardClauseSpacingAllowsSeparatedGuardClauses(t *testing.T) {
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
	fileSet, file := parseGoSource(t, source)

	violations := CheckGuardClauseSpacing(fileSet, file)
	if len(violations) != 0 {
		t.Fatalf("expected no guard-clause violations, got: %#v", violations)
	}
}
