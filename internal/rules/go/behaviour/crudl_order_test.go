package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsCRUDLOrderViolation(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "client", "application", "port", "message")
	sourcePath := filepath.Join(portsDirectory, "repository.go")
	sourceCode := `package ports

type MessageRepository interface {
	Load(value string) (err error)
	Save(value string) (err error)
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`[go/order/crudl] method "Save" in interface "MessageRepository" is out of CRUD-L order`,
	) {
		t.Fatalf("expected CRUD-L order violation, got:\n%s", output)
	}
}

func TestStylecheckPassesValidCRUDLOrder(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "client", "application", "port", "message")
	portsPath := filepath.Join(portsDirectory, "repository.go")
	portsSource := `package ports

type MessageRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`
	writeSourceFile(t, portsPath, portsSource)

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, output:\n%s", output)
	}
}
