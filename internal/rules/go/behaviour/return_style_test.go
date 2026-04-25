package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

/* ---------------------------------------- Return Values --------------------------------------- */

func TestStylecheckReportsUnnamedReturns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Config struct{}

func Bad(value string) error {
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
		`[go/returns/named-values] function "Bad" has unnamed return values`,
	) {
		t.Fatalf("expected unnamed return violation, got:\n%s", output)
	}
}

func TestStylecheckPassesValidFile(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, output:\n%s", output)
	}
}

func TestStylecheckReportsPlaceholderReturnNames(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

func Bad() (result0 string) {
	return "bad"
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`[go/returns/no-placeholder-names] function "Bad" uses placeholder return name "result0"`,
	) {
		t.Fatalf("expected placeholder return-name violation, got:\n%s", output)
	}
}

func TestStylecheckReportsNakedReturns(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`[go/returns/no-naked-returns] function "Bad" uses a naked return`,
	) {
		t.Fatalf("expected naked-return violation, got:\n%s", output)
	}
}
