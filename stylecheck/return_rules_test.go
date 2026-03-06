package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsUnnamedReturns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Config struct{}

func Bad(value string) error {
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

	if !strings.Contains(output, `[2.2] function "Bad" has unnamed return values`) {
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
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
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
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `[2.2] function "Bad" uses placeholder return name "result0"`) {
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
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `[2.2] function "Bad" uses a naked return`) {
		t.Fatalf("expected naked-return violation, got:\n%s", output)
	}
}
