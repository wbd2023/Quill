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

func Worse(secretToken string, value string) (err error) {
	return nil
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

	expectDiagnosticMessage(
		t,
		result,
		`[go/parameters/context-first] ctx must be the first parameter in "Bad"`,
	)

	expectDiagnosticMessage(
		t,
		result,
		`[go/parameters/secrets-last] secret parameters must be last in "Worse"`,
	)
}

func TestGoStyleReportsTypeElision(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

func Bad(left, right string) (err error) {
	return nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected type-elision failure, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(
		t,
		result,
		`[go/parameters/no-type-elision] type elision: parameters left, right share a type`,
	)
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(
		t,
		result,
		`[go/constructors/category-order] repository parameter appears after `+
			`adapter parameter in constructor "NewThing"`,
	)
}

func TestGoStylePassesValidParameterRules(t *testing.T) {
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
	secretToken string,
) (thing *Thing) {
	_ = userRepository
	_ = mailService
	_ = relayClient
	_ = relayURL
	_ = secretToken
	return &Thing{}
}

func Good(ctx Context, value string, secretToken string) (err error) {
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
