package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsParameterOrderViolations(t *testing.T) {
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
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `[2.7] ctx must be the first parameter in "Bad"`) {
		t.Fatalf("expected ctx-order violation, got:\n%s", output)
	}

	if !strings.Contains(output, `[2.7] secret parameters must be last in "Worse"`) {
		t.Fatalf("expected secret-order violation, got:\n%s", output)
	}
}

func TestStylecheckReportsConstructorParameterOrderViolation(t *testing.T) {
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
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`[2.8] repository parameter appears after adapter parameter in constructor "NewThing"`,
	) {
		t.Fatalf("expected constructor-order violation, got:\n%s", output)
	}
}

func TestStylecheckPassesValidParameterAndConstructorOrder(t *testing.T) {
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
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}
