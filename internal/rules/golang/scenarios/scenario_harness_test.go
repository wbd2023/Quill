package scenarios

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rules/golang"
)

/* ------------------------------------------- Harness ------------------------------------------ */

func runGoStyleResult(
	t *testing.T,
	targetDirectory string,
) (result contract.ExecutionResult, err error) {
	t.Helper()

	return runGoStyleResultWithPolicy(t, targetDirectory, profiles.Current(t))
}

func runGoStyleResultWithPolicy(
	t *testing.T,
	targetDirectory string,
	config policy.Config,
) (result contract.ExecutionResult, err error) {
	t.Helper()

	result, err = golang.CheckDirectories(
		targetDirectory,
		[]string{targetDirectory},
		config,
	)
	return result, err
}

func expectDiagnosticMessage(
	t *testing.T,
	result contract.ExecutionResult,
	fragment string,
) {
	t.Helper()

	if hasDiagnosticText(result, fragment) {
		return
	}

	t.Fatalf("expected diagnostic containing %q, got %+v", fragment, result.Diagnostics)
}

func rejectDiagnosticMessage(
	t *testing.T,
	result contract.ExecutionResult,
	fragment string,
) {
	t.Helper()

	if hasDiagnosticText(result, fragment) {
		t.Fatalf("unexpected diagnostic containing %q: %+v", fragment, result.Diagnostics)
	}
}

func hasDiagnosticText(result contract.ExecutionResult, fragment string) (found bool) {
	for _, diagnostic := range result.Diagnostics {
		text := diagnostic.Code + " " +
			diagnostic.File + " " +
			diagnostic.Message + " [" +
			diagnostic.Code + "] " +
			diagnostic.Message
		if strings.Contains(text, fragment) {
			return true
		}
	}

	return false
}

func writeTypeAwareDomainFixture(t *testing.T, rootDirectory string) {
	t.Helper()

	fixtures.WriteFile(t, rootDirectory, "go.mod", "module example\n\ngo 1.24.5\n")
	fixtures.WriteFile(
		t,
		rootDirectory,
		"internal/core/domain/types.go",
		`package domain

type IdentityID string
`,
	)
}

func writeSourceFile(t *testing.T, path string, contents string) {
	t.Helper()

	fixtures.WritePath(t, path, contents)
}
