package checks

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

/* --------------------------------------- Ordering Checks -------------------------------------- */

func TestCheckScannerEntrypointOrderRejectsHelperBeforeCheck(t *testing.T) {
	fileSet, file := parseGoSource(t, `package text

func helper() {}

func CheckThing() {}
`)

	violations := CheckScannerEntrypointOrder(
		fileSet,
		file,
		"/repo/tools/internal/rules/text/example.go",
	)
	if len(violations) != 1 || violations[0].Rule != DiagnosticScannerEntrypointOrder {
		t.Fatalf("expected scanner entrypoint violation, got %#v", violations)
	}
}

func TestCheckScannerEntrypointOrderAcceptsCheckBeforeHelpers(t *testing.T) {
	fileSet, file := parseGoSource(t, `package text

func CheckThing() {}

func helper() {}
`)

	violations := CheckScannerEntrypointOrder(
		fileSet,
		file,
		"/repo/tools/internal/rules/text/example.go",
	)
	if len(violations) != 0 {
		t.Fatalf("expected scanner order to pass, got %#v", violations)
	}
}

func TestCheckTestHygieneRejectsHelperBeforeTests(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

import "testing"

func helper(t *testing.T) {
	t.Helper()
}

func TestExample(t *testing.T) {}
`)

	violations := CheckTestHygiene(fileSet, file, "example_test.go")
	if !hasViolation(violations, DiagnosticTestHelperOrder) {
		t.Fatalf("expected test helper order violation, got %#v", violations)
	}
}

func TestCheckTestHygieneRejectsHelperBetweenTests(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

import "testing"

func TestFirst(t *testing.T) {}

func helper(t *testing.T) {
	t.Helper()
}

func TestSecond(t *testing.T) {}
`)

	violations := CheckTestHygiene(fileSet, file, "example_test.go")
	if !hasViolation(violations, DiagnosticTestHelperOrder) {
		t.Fatalf("expected test helper order violation, got %#v", violations)
	}
}

func TestCheckTestHygieneAllowsHelpersAfterTests(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

import "testing"

func TestExample(t *testing.T) {}

func helper(t *testing.T) {
	t.Helper()
}
`)

	violations := CheckTestHygiene(fileSet, file, "example_test.go")
	if hasViolation(violations, DiagnosticTestHelperOrder) {
		t.Fatalf("expected helper order to pass, got %#v", violations)
	}
}

/* ---------------------------------------- Parse Helpers --------------------------------------- */

func parseGoSource(t *testing.T, source string) (fileSet *token.FileSet, file *ast.File) {
	t.Helper()

	fileSet = token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "example.go", source, parser.ParseComments)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	return fileSet, file
}

func hasViolation(violations []Violation, rule string) (found bool) {
	for _, violation := range violations {
		if violation.Rule == rule {
			return true
		}
	}

	return false
}
