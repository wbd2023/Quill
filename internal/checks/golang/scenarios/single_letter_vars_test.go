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

func AlsoBad(i string) (err error) {
	k := i
	_ = k
	return nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/single-letter-names] single-letter parameter "x" in function "Bad"`,
	)

	expectDiagnosticMessage(t, result, `[go/naming/single-letter-names] single-letter variable "y"`)

	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/single-letter-names] single-letter parameter "i" in function "AlsoBad"`,
	)

	expectDiagnosticMessage(t, result, `[go/naming/single-letter-names] single-letter variable "k"`)
}

func TestGoStylePassesAllowedSingleLetterNames(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

func Good(values []int) (err error) {
	for i := range values {
		_ = i
	}

	for j := 0; j < len(values); j++ {
		_ = values[j]
	}

	for var k = 0; k < len(values); k++ {
		_ = values[k]
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
