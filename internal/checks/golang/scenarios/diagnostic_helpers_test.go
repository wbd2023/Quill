package scenarios

import (
	"path/filepath"
	"strings"
	"testing"

	"ciphera/tools/internal/style"
)

type diagnosticMatch struct {
	code    string
	file    string
	line    int
	column  int
	message string
}

func expectDiagnosticMessage(
	t *testing.T,
	result style.ExecutionResult,
	fragment string,
) {
	t.Helper()

	expectDiagnostic(t, result, diagnosticMatchFromFragment(fragment))
}

func rejectDiagnosticMessage(
	t *testing.T,
	result style.ExecutionResult,
	fragment string,
) {
	t.Helper()

	if hasDiagnostic(result, diagnosticMatchFromFragment(fragment)) {
		t.Fatalf("unexpected diagnostic containing %q: %+v", fragment, result.Diagnostics)
	}
}

func expectDiagnostic(t *testing.T, result style.ExecutionResult, expected diagnosticMatch) {
	t.Helper()

	if hasDiagnostic(result, expected) {
		return
	}

	t.Fatalf("expected diagnostic %+v, got %+v", expected, result.Diagnostics)
}

func hasDiagnostic(result style.ExecutionResult, expected diagnosticMatch) (found bool) {
	for _, diagnostic := range result.Diagnostics {
		if diagnosticMatches(diagnostic, expected) {
			return true
		}
	}

	return false
}

func diagnosticMatches(diagnostic style.Diagnostic, expected diagnosticMatch) (matches bool) {
	if expected.code != "" && diagnostic.Code != expected.code {
		return false
	}

	if expected.file != "" && filepath.ToSlash(diagnostic.File) != filepath.ToSlash(expected.file) {
		return false
	}

	if expected.line != 0 && diagnostic.Line != expected.line {
		return false
	}

	if expected.column != 0 && diagnostic.Column != expected.column {
		return false
	}

	if expected.message != "" && !strings.Contains(diagnostic.Message, expected.message) {
		return false
	}

	return true
}

func diagnosticMatchFromFragment(fragment string) (expected diagnosticMatch) {
	if code, message, ok := strings.Cut(strings.TrimPrefix(fragment, "["), "] "); ok &&
		strings.HasPrefix(fragment, "[") {
		return diagnosticMatch{
			code:    code,
			message: message,
		}
	}

	return diagnosticMatch{message: fragment}
}
