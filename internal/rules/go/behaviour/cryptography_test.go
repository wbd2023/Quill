package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsCryptographyImportViolations(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected cryptography failure, output:\n%s", output)
	}

	if !strings.Contains(output, "production code must not import math/rand") {
		t.Fatalf("expected math/rand violation, got:\n%s", output)
	}

	if !strings.Contains(output, "deprecated cryptographic package crypto/sha1") {
		t.Fatalf("expected deprecated crypto violation, got:\n%s", output)
	}
}

func TestStylecheckAcceptsModernCryptographyImports(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "core", "crypto", "good.go")
	sourceCode := `package crypto

import "crypto/rand"

var _ = rand.Reader
`
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected cryptography fixture to pass, output:\n%s", output)
	}
}
