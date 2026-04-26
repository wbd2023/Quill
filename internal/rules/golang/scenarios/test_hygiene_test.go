package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsTestHygieneViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(
		tempDir,
		"internal",
		"client",
		"application",
		"service",
		"sample_test.go",
	)
	sourceCode := `package service

import (
	"os"
	"testing"
	"time"
)

func writeFixture(t *testing.T) {
	_, _ = os.MkdirTemp("", "fixture")
	os.Setenv("HOME", "/tmp/example")
	time.Sleep(time.Second)
}

func TestSomething(t *testing.T) {}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected test hygiene failure, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(
		t,
		result,
		"test helpers that accept testing handles must call Helper()",
	)

	expectDiagnosticMessage(t, result, "tests must use t.Setenv() instead of os.Setenv()")

	expectDiagnosticMessage(t, result, "tests must use t.TempDir() instead of os.MkdirTemp()")

	expectDiagnosticMessage(
		t,
		result,
		"tests must avoid time.Sleep() when a deterministic signal is possible",
	)
}

func TestGoStyleAcceptsTestHygienePatterns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(
		tempDir,
		"internal",
		"client",
		"application",
		"service",
		"sample_test.go",
	)
	sourceCode := `package service

import "testing"

func TestSomething(t *testing.T) {
	directory := t.TempDir()
	t.Setenv("HOME", directory)
}

func writeFixture(t *testing.T) {
	t.Helper()
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf(
			"expected valid test hygiene fixture to pass, diagnostics: %#v",
			result.Diagnostics,
		)
	}
}
