package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsCryptographyImportViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "core", "crypto", "bad.go")
	sourceCode := `package crypto

import (
	"crypto/sha1"
	"math/rand"
)

var _ = sha1.New
var _ = rand.Int
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected cryptography failure, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, "production code must not import math/rand")

	expectDiagnosticMessage(t, result, "deprecated cryptographic package crypto/sha1")
}

func TestGoStyleAcceptsModernCryptographyImports(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "core", "crypto", "good.go")
	sourceCode := `package crypto

import "crypto/rand"

var _ = rand.Reader
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected cryptography fixture to pass, diagnostics: %#v", result.Diagnostics)
	}
}
