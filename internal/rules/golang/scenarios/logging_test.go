package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsStructuredLoggingViolations(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected structured logging failure, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "structured log key \"Path\" must be lower-case ASCII") {
		t.Fatalf("expected structured log key violation, got: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "structured log calls must use key/value pairs") {
		t.Fatalf("expected structured log key/value violation, got: %#v", result.Diagnostics)
	}
}

func TestGoStyleAcceptsStructuredLoggingWithStableKeys(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf(
			"expected structured logging fixture to pass, diagnostics: %#v",
			result.Diagnostics,
		)
	}
}

func TestGoStyleRejectsSecretStructuredLoggingFields(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected structured logging failure, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "structured logs must not include secrets") {
		t.Fatalf("expected secret logging violation, got: %#v", result.Diagnostics)
	}
}
