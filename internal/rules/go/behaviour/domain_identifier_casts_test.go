package behaviour

import (
	"path/filepath"
	"strings"
	"testing"

	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/profile"
)

/* ----------------------------------- Domain Identifier Casts ---------------------------------- */

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
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
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
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, output:\n%s", output)
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
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `direct cast to domain.IdentityID is disallowed`) {
		t.Fatalf("expected type-aware direct-cast violation, got:\n%s", output)
	}
}

func TestStylecheckUsesProfileDomainIdentifierVocabulary(t *testing.T) {
	tempDir := t.TempDir()
	policy := profiles.Current(t)
	policy.Naming.GoDomainIdentifiers = profile.GoDomainIdentifierConfig{
		"SessionKey": {"ParseSessionKey"},
	}

	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "project/internal/core/domain"

func Bad(raw string) (id domain.SessionKey, err error) {
	id = domain.SessionKey(raw)
	return id, nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheckWithPolicy(t, tempDir, policy)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `use ParseSessionKey`) {
		t.Fatalf("expected alternate constructor in output, got:\n%s", output)
	}
}
