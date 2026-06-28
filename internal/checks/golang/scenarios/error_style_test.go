package scenarios

import (
	"path/filepath"
	"testing"
)

/* --------------------------------------- Error Handling --------------------------------------- */

func TestGoStyleReportsErrorHandlingViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(
		tempDir,
		"internal",
		"client",
		"application",
		"usecase",
		"message",
		"sample.go",
	)
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
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, "error context must be lowercase (errors.New)")

	expectDiagnosticMessage(t, result, "error context must not end with punctuation (errors.New)")

	expectDiagnosticMessage(
		t,
		result,
		"error context must not include secrets in fmt.Errorf arguments",
	)
}

func TestGoStyleReportsSentinelErrorsOutsideDomainErrorsFile(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(
		tempDir,
		"internal",
		"client",
		"application",
		"usecase",
		"message",
		"sample.go",
	)
	sourceCode := `package services

import "errors"

var ErrServiceFailed = errors.New("service failed")
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
		"sentinel errors must be declared in internal/core/domain/errors.go",
	)
}

func TestGoStyleReportsAdapterBareErrReturn(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(
		tempDir,
		"internal",
		"client",
		"adapters",
		"outbound",
		"persistence",
		"filestore",
		"sample.go",
	)
	sourceCode := `package storage

import "errors"

func BadAdapter() (value string, err error) {
	err = errors.New("load failed")
	return value, err
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
		"adapter error returns must wrap low-level errors with context (%w)",
	)
}

func TestGoStylePassesAdapterWrappedErrorReturn(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(
		tempDir,
		"internal",
		"client",
		"adapters",
		"outbound",
		"persistence",
		"filestore",
		"sample.go",
	)
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
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, diagnostics: %#v", result.Diagnostics)
	}
}
