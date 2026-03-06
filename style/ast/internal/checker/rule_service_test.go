package checker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsServiceTypeNamingViolation(t *testing.T) {
	tempDir := t.TempDir()
	serviceDirectory := filepath.Join(tempDir, "internal", "core", "services", "message")
	if err := os.MkdirAll(serviceDirectory, 0o700); err != nil {
		t.Fatalf("mkdir services: %v", err)
	}

	sourcePath := filepath.Join(serviceDirectory, "service.go")
	sourceCode := `package message

type Manager struct{}
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
		`[2.2] exported type "Manager" should end with Service, UseCase, or Config`,
	) {
		t.Fatalf("expected service-type naming violation, got:\n%s", output)
	}
}

func TestStylecheckPassesValidServiceTypeNaming(t *testing.T) {
	tempDir := t.TempDir()
	serviceDirectory := filepath.Join(tempDir, "internal", "core", "services", "message")
	if err := os.MkdirAll(serviceDirectory, 0o700); err != nil {
		t.Fatalf("mkdir services: %v", err)
	}

	servicePath := filepath.Join(serviceDirectory, "service.go")
	serviceSource := `package message

type MessageService struct{}
`
	if err := os.WriteFile(servicePath, []byte(serviceSource), 0o600); err != nil {
		t.Fatalf("write service source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}
