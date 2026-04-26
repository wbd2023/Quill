package scenarios

import (
	"path/filepath"
	"testing"
)

/* ----------------------------------------- Parameters ----------------------------------------- */

func TestGoStyleReportsParameterOrderViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Context struct{}

func Bad(value string, ctx Context) (err error) {
	return nil
}

func Worse(token string, value string) (err error) {
	return nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(
		result,
		`[go/parameters/context-first] ctx must be the first parameter in "Bad"`,
	) {
		t.Fatalf("expected ctx-order violation, got: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(
		result,
		`[go/parameters/secrets-last] secret parameters must be last in "Worse"`,
	) {
		t.Fatalf("expected secret-order violation, got: %#v", result.Diagnostics)
	}
}

func TestGoStyleReportsConstructorParameterOrderViolation(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type UserRepository interface{}
type MailService interface{}
type RelayClient struct{}
type Thing struct{}

func NewThing(relayClient *RelayClient, userRepository UserRepository) (thing *Thing) {
	_ = relayClient
	_ = userRepository
	return &Thing{}
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	if !hasDiagnosticText(
		result,
		`[go/constructors/category-order] repository parameter appears after `+
			`adapter parameter in constructor "NewThing"`,
	) {
		t.Fatalf("expected constructor-order violation, got: %#v", result.Diagnostics)
	}
}

func TestGoStylePassesValidParameterAndConstructorOrder(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Context struct{}
type UserRepository interface{}
type MailService interface{}
type RelayClient struct{}
type Thing struct{}

func NewThing(
	userRepository UserRepository,
	mailService MailService,
	relayClient *RelayClient,
	relayURL string,
	token string,
) (thing *Thing) {
	_ = userRepository
	_ = mailService
	_ = relayClient
	_ = relayURL
	_ = token
	return &Thing{}
}

func Good(ctx Context, value string, token string) (err error) {
	_ = value
	return nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, diagnostics: %#v", result.Diagnostics)
	}
}
