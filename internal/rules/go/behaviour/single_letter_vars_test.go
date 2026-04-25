package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsSingleLetterNames(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`[go/naming/single-letter-names] single-letter parameter "x" in function "Bad"`,
	) {
		t.Fatalf("expected single-letter parameter violation, got:\n%s", output)
	}

	if !strings.Contains(output, `[go/naming/single-letter-names] single-letter variable "y"`) {
		t.Fatalf("expected single-letter variable violation, got:\n%s", output)
	}
}

func TestStylecheckPassesAllowedSingleLetterNames(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, output:\n%s", output)
	}
}
