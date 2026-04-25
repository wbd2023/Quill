package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsShellInterpolationProcessExecution(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected process execution failure, output:\n%s", output)
	}

	if !strings.Contains(output, "process execution must avoid shell interpolation") {
		t.Fatalf("expected shell interpolation violation, got:\n%s", output)
	}
}

func TestStylecheckAcceptsDirectArgumentProcessExecution(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected direct exec fixture to pass, output:\n%s", output)
	}
}
