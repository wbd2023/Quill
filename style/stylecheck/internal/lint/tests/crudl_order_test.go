package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsCRUDLOrderViolation(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "client", "application", "port", "message")
	if err := os.MkdirAll(portsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir ports: %v", err)
	}

	sourcePath := filepath.Join(portsDirectory, "repository.go")
	sourceCode := `package ports

type MessageRepository interface {
	Load(value string) (err error)
	Save(value string) (err error)
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
		`[2.5] method "Save" in interface "MessageRepository" is out of CRUD-L order`,
	) {
		t.Fatalf("expected CRUD-L order violation, got:\n%s", output)
	}
}

func TestStylecheckPassesValidCRUDLOrder(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "client", "application", "port", "message")
	if err := os.MkdirAll(portsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir ports: %v", err)
	}

	portsPath := filepath.Join(portsDirectory, "repository.go")
	portsSource := `package ports

type MessageRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`
	if err := os.WriteFile(portsPath, []byte(portsSource), 0o600); err != nil {
		t.Fatalf("write ports source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}
