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

	if !hasDiagnosticText(result, "test helpers that accept testing handles must call Helper()") {
		t.Fatalf("expected missing Helper violation, got: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "tests must use t.Setenv() instead of os.Setenv()") {
		t.Fatalf("expected Setenv violation, got: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "tests must use t.TempDir() instead of os.MkdirTemp()") {
		t.Fatalf("expected TempDir violation, got: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(
		result,
		"tests must avoid time.Sleep() when a deterministic signal is possible",
	) {
		t.Fatalf("expected time.Sleep violation, got: %#v", result.Diagnostics)
	}
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

func writeFixture(t *testing.T) {
	t.Helper()
}

func TestSomething(t *testing.T) {
	directory := t.TempDir()
	t.Setenv("HOME", directory)
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
