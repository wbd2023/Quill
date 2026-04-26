package scenarios

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/rules/golang"
)

/* --------------------------------------- Rule Selection --------------------------------------- */

func TestGoStyleCheckRunsOnlyRequestedDiagnosticFamily(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "relay", "bootstrap", "sample.go")
	sourceCode := `package bootstrap

import (
	"log/slog"
	"project/internal/core/domain"
)

func Bad(raw string) (id domain.IdentityID, err error) {
	slog.Info("access", "Path", "/")
	id = domain.IdentityID(raw)
	return id, nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runSelectedGoStyleCheck(t, tempDir, golang.GoCheckLogging)
	if err == nil {
		t.Fatalf("expected logging check to fail, result: %+v", result)
	}

	expectDiagnosticMessage(t, result, "structured log key \"Path\" must be lower-case ASCII")
	rejectDiagnosticMessage(t, result, "direct cast to domain.IdentityID")

	result, err = runSelectedGoStyleCheck(t, tempDir, golang.GoCheckDomainIdentifiers)
	if err == nil {
		t.Fatalf("expected domain identifier check to fail, result: %+v", result)
	}

	expectDiagnosticMessage(t, result, "direct cast to domain.IdentityID")
	rejectDiagnosticMessage(t, result, "structured log key")
}

func TestGoStyleExecutorReportsGuardSpacingRule(t *testing.T) {
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

	result, err := runSelectedGoStyleCheck(t, tempDir, golang.GoCheckGuardClauseSpacing)
	if err == nil {
		t.Fatalf("expected guard-clause spacing check to fail, result: %+v", result)
	}

	expectDiagnosticMessage(
		t,
		result,
		"consecutive guard clauses should be separated by a blank line",
	)
}

func TestGoStyleExecutorReportsSwitchSpacingRule(t *testing.T) {
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

	result, err := runSelectedGoStyleCheck(t, tempDir, golang.GoCheckSwitchCaseSpacing)
	if err == nil {
		t.Fatalf("expected switch-case spacing check to fail, result: %+v", result)
	}

	expectDiagnosticMessage(t, result, "non-trivial switch statements should separate case blocks")
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func runSelectedGoStyleCheck(
	t *testing.T,
	targetDirectory string,
	checkName string,
) (result contract.ExecutionResult, err error) {
	t.Helper()

	result, err = golang.CheckDirectories(
		targetDirectory,
		[]string{filepath.Join(targetDirectory, "internal")},
		profiles.Current(t),
		checkName,
	)
	return result, err
}
