package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsStructuredLoggingViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "relay", "bootstrap", "logging.go")
	sourceCode := `package bootstrap

import "log/slog"

func BadLogging() {
	logger := slog.Default()
	logger.Info("access", "Path", "/")
	slog.Warn("access", "ip")
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected structured logging failure, output:\n%s", output)
	}

	if !strings.Contains(output, "structured log key \"Path\" must be lower-case ASCII") {
		t.Fatalf("expected structured log key violation, got:\n%s", output)
	}

	if !strings.Contains(output, "structured log calls must use key/value pairs") {
		t.Fatalf("expected structured log key/value violation, got:\n%s", output)
	}
}

func TestStylecheckAcceptsStructuredLoggingWithStableKeys(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "relay", "bootstrap", "logging.go")
	sourceCode := `package bootstrap

import "log/slog"

func GoodLogging() {
	logger := slog.Default()
	logger.Info("access", "ip", "127.0.0.1", "path", "/")
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected structured logging fixture to pass, output:\n%s", output)
	}
}

func TestStylecheckRejectsSecretStructuredLoggingFields(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "relay", "bootstrap", "logging.go")
	sourceCode := `package bootstrap

import "log/slog"

func BadLogging(passphrase string) {
	logger := slog.Default()
	logger.Info("access", "passphrase", passphrase)
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected structured logging failure, output:\n%s", output)
	}

	if !strings.Contains(output, "structured logs must not include secrets") {
		t.Fatalf("expected secret logging violation, got:\n%s", output)
	}
}
