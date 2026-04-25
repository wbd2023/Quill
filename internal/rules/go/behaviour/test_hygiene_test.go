package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsTestHygieneViolations(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected test hygiene failure, output:\n%s", output)
	}

	if !strings.Contains(output, "test helpers that accept testing handles must call Helper()") {
		t.Fatalf("expected missing Helper violation, got:\n%s", output)
	}

	if !strings.Contains(output, "tests must use t.Setenv() instead of os.Setenv()") {
		t.Fatalf("expected Setenv violation, got:\n%s", output)
	}

	if !strings.Contains(output, "tests must use t.TempDir() instead of os.MkdirTemp()") {
		t.Fatalf("expected TempDir violation, got:\n%s", output)
	}

	if !strings.Contains(
		output,
		"tests must avoid time.Sleep() when a deterministic signal is possible",
	) {
		t.Fatalf("expected time.Sleep violation, got:\n%s", output)
	}
}

func TestStylecheckAcceptsTestHygienePatterns(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected valid test hygiene fixture to pass, output:\n%s", output)
	}
}
