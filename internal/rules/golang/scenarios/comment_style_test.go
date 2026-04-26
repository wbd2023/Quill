package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsInlineCommentStyleViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "core", "services", "sample.go")
	sourceCode := `package services

func BadInline(value string) (err error) {
	_ = value // Uppercase
	_ = value // lower-case.
	return nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "inline trailing comment should start lower-case") {
		t.Fatalf("expected inline-case violation, got: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "inline trailing comment should not end with punctuation") {
		t.Fatalf("expected inline-punctuation violation, got: %#v", result.Diagnostics)
	}
}
