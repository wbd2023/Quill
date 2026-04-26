package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsSingleLetterNames(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

func Bad(x string) (err error) {
	y := x
	_ = y
	return nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(
		result,
		`[go/naming/single-letter-names] single-letter parameter "x" in function "Bad"`,
	) {
		t.Fatalf("expected single-letter parameter violation, got: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, `[go/naming/single-letter-names] single-letter variable "y"`) {
		t.Fatalf("expected single-letter variable violation, got: %#v", result.Diagnostics)
	}
}

func TestGoStylePassesAllowedSingleLetterNames(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

func Good(values []int) (err error) {
	for i := range values {
		_ = i
	}

	return nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, diagnostics: %#v", result.Diagnostics)
	}
}
