package scenarios

import (
	"path/filepath"
	"testing"

	"github.com/wbd2023/Quill/internal/checks/gopolicy"
)

/* ----------------------------------- Domain Identifier Casts ---------------------------------- */

func TestGoStyleReportsDirectDomainValueCast(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "profile/internal/core/domain"

func Bad(raw string) (id domain.IdentityID, err error) {
	id = domain.IdentityID(raw)
	return id, nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, `direct cast to domain.IdentityID is disallowed`)
}

func TestGoStylePassesDomainValueParserUsage(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "profile/internal/core/domain"

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

func TestGoStyleReportsTypeAliasToDomainValueCast(t *testing.T) {
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, `direct cast to domain.IdentityID is disallowed`)
}

func TestGoStyleUsesProfileDomainValueVocabulary(t *testing.T) {
	tempDir := t.TempDir()
	config := scenarioConfig(t)
	updateGoConfigForTest(t, &config, func(goConfig *gopolicy.Config) {
		goConfig.DomainValues.RequiredConstructors = gopolicy.DomainValueConstructors{
			"SessionKey": {"ParseSessionKey"},
		}
	})

	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "profile/internal/core/domain"

func Bad(raw string) (id domain.SessionKey, err error) {
	id = domain.SessionKey(raw)
	return id, nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResultWithPolicy(t, tempDir, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, `use ParseSessionKey`)
}
