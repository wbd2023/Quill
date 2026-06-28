package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsCRUDLOrderViolation(t *testing.T) {
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
		`[go/order/crudl] method "Save" in interface "MessageRepository" is out of CRUD-L order`,
	)
}

func TestGoStylePassesValidCRUDLOrder(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, diagnostics: %#v", result.Diagnostics)
	}
}
