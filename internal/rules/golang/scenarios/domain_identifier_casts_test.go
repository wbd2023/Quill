package scenarios

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/policy"
)

/* ----------------------------------- Domain Identifier Casts ---------------------------------- */

func TestGoStyleReportsDirectDomainIdentifierCast(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, `direct cast to domain.IdentityID is disallowed`)
}

func TestGoStylePassesDomainIdentifierParserUsage(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, diagnostics: %#v", result.Diagnostics)
	}
}

func TestGoStyleReportsTypeAliasToDomainIdentifierCast(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, `direct cast to domain.IdentityID is disallowed`)
}

func TestGoStyleUsesProfileDomainIdentifierVocabulary(t *testing.T) {
	tempDir := t.TempDir()
	config := profiles.Current(t)
	config.Go.DomainIdentifierConstructors = policy.GoDomainIdentifierConstructors{
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

	result, err := runGoStyleResultWithPolicy(t, tempDir, config)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, `use ParseSessionKey`)
}
