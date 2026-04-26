package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsShellInterpolationProcessExecution(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "exec.go")
	sourceCode := `package service

import "os/exec"

func BadExec(commandText string) (command *exec.Cmd) {
	command = exec.Command("sh", "-c", commandText)
	return command
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected process execution failure, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(t, result, "process execution must avoid shell interpolation")
}

func TestGoStyleAcceptsDirectArgumentProcessExecution(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "exec.go")
	sourceCode := `package service

import "os/exec"

func GoodExec() (command *exec.Cmd) {
	command = exec.Command("go", "version")
	return command
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected direct exec fixture to pass, diagnostics: %#v", result.Diagnostics)
	}
}
