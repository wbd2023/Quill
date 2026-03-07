package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsInlineCommentStyleViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "core", "services", "sample.go")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
		t.Fatalf("mkdir sample directory: %v", err)
	}

	sourceCode := `package services

func BadInline(value string) (err error) {
	_ = value // Uppercase
	_ = value // lower-case.
	return nil
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "inline trailing comment should start lower-case") {
		t.Fatalf("expected inline-case violation, got:\n%s", output)
	}

	if !strings.Contains(output, "inline trailing comment should not end with punctuation") {
		t.Fatalf("expected inline-punctuation violation, got:\n%s", output)
	}
}
