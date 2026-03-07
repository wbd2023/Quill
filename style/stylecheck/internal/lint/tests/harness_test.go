package lint

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

/* ------------------------------------------- Helpers ------------------------------------------ */

func runStylecheck(targetDirectory string) (output string, err error) {
	command := exec.Command("go", "run", ".", targetDirectory)
	command.Dir = stylecheckModuleDirectory()

	rawOutput, runErr := command.CombinedOutput()
	return string(rawOutput), runErr
}

func stylecheckModuleDirectory() (directory string) {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".."))
}

func writeTypeAwareDomainFixture(t *testing.T, rootDirectory string) {
	t.Helper()

	goModPath := filepath.Join(rootDirectory, "go.mod")
	goModContent := "module example\n\ngo 1.24.5\n"
	if err := os.WriteFile(goModPath, []byte(goModContent), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	domainDirectory := filepath.Join(rootDirectory, "internal", "core", "domain")
	if err := os.MkdirAll(domainDirectory, 0o700); err != nil {
		t.Fatalf("mkdir domain fixture: %v", err)
	}

	domainTypesPath := filepath.Join(domainDirectory, "types.go")
	domainTypesSource := `package domain

type IdentityID string
`
	if err := os.WriteFile(domainTypesPath, []byte(domainTypesSource), 0o600); err != nil {
		t.Fatalf("write domain fixture types: %v", err)
	}
}
