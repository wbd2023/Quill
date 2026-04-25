package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsInlineCommentStyleViolations(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "inline trailing comment should start lower-case") {
		t.Fatalf("expected inline-case violation, got:\n%s", output)
	}

	if !strings.Contains(output, "inline trailing comment should not end with punctuation") {
		t.Fatalf("expected inline-punctuation violation, got:\n%s", output)
	}
}
