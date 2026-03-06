package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsErrorHandlingViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "core", "services", "sample.go")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
		t.Fatalf("mkdir sample directory: %v", err)
	}

	sourceCode := `package services

import (
	"errors"
	"fmt"
)

func BadErrorStyle(secretToken string) (err error) {
	_ = errors.New("Bad.")
	_ = fmt.Errorf("failed auth for %s", secretToken)
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

	if !strings.Contains(output, "error context must be lowercase (errors.New)") {
		t.Fatalf("expected lower-case error-context violation, got:\n%s", output)
	}

	if !strings.Contains(output, "error context must not end with punctuation (errors.New)") {
		t.Fatalf("expected punctuation error-context violation, got:\n%s", output)
	}

	if !strings.Contains(output, "error context must not include secrets in fmt.Errorf arguments") {
		t.Fatalf("expected secret-argument violation, got:\n%s", output)
	}
}

func TestStylecheckReportsSentinelErrorsOutsideDomainErrorsFile(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "core", "services", "sample.go")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
		t.Fatalf("mkdir sample directory: %v", err)
	}

	sourceCode := `package services

import "errors"

var ErrServiceFailed = errors.New("service failed")
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		"sentinel errors must be declared in internal/core/domain/errors.go",
	) {
		t.Fatalf("expected sentinel-location violation, got:\n%s", output)
	}
}

func TestStylecheckReportsAdapterBareErrReturn(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "adapters", "storage", "sample.go")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
		t.Fatalf("mkdir adapter directory: %v", err)
	}

	sourceCode := `package storage

import "errors"

func BadAdapter() (value string, err error) {
	err = errors.New("load failed")
	return value, err
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		"adapter error returns must wrap low-level errors with context (%w)",
	) {
		t.Fatalf("expected adapter-wrap violation, got:\n%s", output)
	}
}

func TestStylecheckPassesAdapterWrappedErrorReturn(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "adapters", "storage", "sample.go")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
		t.Fatalf("mkdir adapter directory: %v", err)
	}

	sourceCode := `package storage

import "fmt"

func GoodAdapter() (value string, err error) {
	err = fmt.Errorf("load failed")
	if err != nil {
		return "", fmt.Errorf("load adapter: %w", err)
	}
	return value, nil
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
