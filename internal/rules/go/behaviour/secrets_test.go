package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsHardCodedSecretLiterals(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected secret literal failure, output:\n%s", output)
	}

	if !strings.Contains(output, "source code must not hard-code secret-like string literals") {
		t.Fatalf("expected secret literal violation, got:\n%s", output)
	}
}

func TestStylecheckIgnoresNonSensitiveLiterals(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected non-sensitive fixture to pass, output:\n%s", output)
	}
}
