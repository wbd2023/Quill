package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsHardCodedSecretLiterals(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "bad.go")
	sourceCode := `package service

const bootstrapToken = "A1B2C3D4E5F6G7H8"

func Token() (token string) {
	token = bootstrapToken
	return token
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected secret literal failure, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "source code must not hard-code secret-like string literals") {
		t.Fatalf("expected secret literal violation, got: %#v", result.Diagnostics)
	}
}

func TestGoStyleIgnoresNonSensitiveLiterals(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "good.go")
	sourceCode := `package service

const relayURL = "https://example.com"

func RelayURL() (value string) {
	value = relayURL
	return value
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected non-sensitive fixture to pass, diagnostics: %#v", result.Diagnostics)
	}
}
