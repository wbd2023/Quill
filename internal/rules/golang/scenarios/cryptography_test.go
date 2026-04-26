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

	if !hasDiagnosticText(result, "production code must not import math/rand") {
		t.Fatalf("expected math/rand violation, got: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(result, "deprecated cryptographic package crypto/sha1") {
		t.Fatalf("expected deprecated crypto violation, got: %#v", result.Diagnostics)
	}
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
