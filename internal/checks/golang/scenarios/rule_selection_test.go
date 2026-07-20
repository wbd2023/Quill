package scenarios

import (
	"path/filepath"
	"testing"

	"github.com/wbd2023/Quill/internal/checks/golang"
	"github.com/wbd2023/Quill/internal/checks/golang/check"
	"github.com/wbd2023/Quill/internal/style"
)

/* --------------------------------------- Rule Selection --------------------------------------- */

func TestGoStyleCheckRunsOnlyRequestedDiagnosticFamily(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "relay", "bootstrap", "sample.go")
	sourceCode := `package bootstrap

import (
	"log/slog"
	"profile/internal/core/domain"
)

func Bad(raw string) (id domain.IdentityID, err error) {
	slog.Info("access", "Path", "/")
	id = domain.IdentityID(raw)
	return id, nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runSelectedGoStyleCheck(t, tempDir, check.Logging)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected logging check to fail, result: %+v", result)
	}

	expectDiagnosticMessage(t, result, "structured log key \"Path\" must be lower-case ASCII")
	rejectDiagnosticMessage(t, result, "direct cast to domain.IdentityID")

	result, err = runSelectedGoStyleCheck(t, tempDir, check.DomainValues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected domain value check to fail, result: %+v", result)
	}

	expectDiagnosticMessage(t, result, "direct cast to domain.IdentityID")
	rejectDiagnosticMessage(t, result, "structured log key")
}

func TestGoStyleDriverReportsGuardSpacingRule(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "example", "guards.go")
	sourceCode := `package example

func Validate(a int, b int) (err error) {
	if a == 0 {
		return nil
	}
	if b == 0 {
		return nil
	}

	return nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runSelectedGoStyleCheck(t, tempDir, check.GuardClauseSpacing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected guard-clause spacing check to fail, result: %+v", result)
	}

	expectDiagnosticMessage(
		t,
		result,
		"consecutive guard clauses should be separated by a blank line",
	)
}

func TestGoStyleDriverReportsSwitchSpacingRule(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "example", "switches.go")
	sourceCode := `package example

func Render(value string) (rendered string) {
	switch value {
	case "a":
		return "A"
	case "b":
		return "B"
	case "c":
		return "C"
	default:
		return "?"
	}
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runSelectedGoStyleCheck(t, tempDir, check.SwitchCaseSpacing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected switch-case spacing check to fail, result: %+v", result)
	}

	expectDiagnosticMessage(t, result, "non-trivial switch statements should separate case blocks")
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func runSelectedGoStyleCheck(
	t *testing.T,
	targetDirectory string,
	checkName string,
) (result style.ExecutionResult, err error) {
	t.Helper()

	config := scenarioConfig(t)
	result, err = golang.CheckDirectories(
		targetDirectory,
		[]string{filepath.Join(targetDirectory, "internal")},
		config.Repository,
		config.PathRoles,
		goConfigForTest(t, config),
		checkName,
	)
	return result, err
}
