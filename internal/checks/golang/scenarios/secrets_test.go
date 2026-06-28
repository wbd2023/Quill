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

	result, _ := runGoStyleResult(t, tempDir)
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected secret literal failure, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, "source code must not hard-code secret-like string literals")
}

func TestGoStyleReportsHardCodedSecretAssignments(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "bad.go")
	sourceCode := `package service

type Client struct {
	Token string
}

func Configure(client *Client) {
	client.Token = "A1B2C3D4E5F6G7H8"
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, _ := runGoStyleResult(t, tempDir)
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected secret assignment failure, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, "source code must not hard-code secret-like string literals")
}

func TestGoStyleReportsHardCodedSecretStructFields(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "bad.go")
	sourceCode := `package service

type Client struct {
	Token string
}

var defaultClient = Client{
	Token: "A1B2C3D4E5F6G7H8",
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, _ := runGoStyleResult(t, tempDir)
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected secret struct-field failure, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, "source code must not hard-code secret-like string literals")
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
