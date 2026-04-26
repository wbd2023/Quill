package scenarios

import (
	"path/filepath"
	"testing"
)

/* ---------------------------------------- Return Values --------------------------------------- */

func TestGoStyleReportsUnnamedReturns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Config struct{}

func Bad(value string) error {
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
		`[go/returns/named-values] function "Bad" has unnamed return values`,
	) {
		t.Fatalf("expected unnamed return violation, got: %#v", result.Diagnostics)
	}
}

func TestGoStylePassesValidFile(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Config struct{}

func Good(value string) (err error) {
	if value == "" {
		return nil
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

func TestGoStyleReportsPlaceholderReturnNames(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

func Bad() (result0 string) {
	return "bad"
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(
		result,
		`[go/returns/no-placeholder-names] function "Bad" uses placeholder return name "result0"`,
	) {
		t.Fatalf("expected placeholder return-name violation, got: %#v", result.Diagnostics)
	}
}

func TestGoStyleReportsNakedReturns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

func Bad(value string) (err error) {
	if value == "" {
		return nil
	}

	return
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(
		result,
		`[go/returns/no-naked-returns] function "Bad" uses a naked return`,
	) {
		t.Fatalf("expected naked-return violation, got: %#v", result.Diagnostics)
	}
}
