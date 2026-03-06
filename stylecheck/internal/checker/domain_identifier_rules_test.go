package checker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsDirectDomainIdentifierCast(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "project/internal/core/domain"

func Bad(raw string) (id domain.IdentityID, err error) {
	id = domain.IdentityID(raw)
	return id, nil
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `direct cast to domain.IdentityID is disallowed`) {
		t.Fatalf("expected direct-cast violation, got:\n%s", output)
	}
}

func TestStylecheckPassesDomainIdentifierParserUsage(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "project/internal/core/domain"

func Good(raw string) (id domain.IdentityID, err error) {
	id, err = domain.ParseIdentityID(raw)
	if err != nil {
		return "", err
	}
	return id, nil
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}

func TestStylecheckReportsTypeAliasToDomainIdentifierCast(t *testing.T) {
	tempDir := t.TempDir()
	writeTypeAwareDomainFixture(t, tempDir)

	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import coredomain "example/internal/core/domain"

type AppIdentityID = coredomain.IdentityID

func Bad(raw string) (id AppIdentityID, err error) {
	id = AppIdentityID(raw)
	return id, nil
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `direct cast to domain.IdentityID is disallowed`) {
		t.Fatalf("expected type-aware direct-cast violation, got:\n%s", output)
	}
}
