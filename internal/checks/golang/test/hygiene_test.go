package test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

func TestCheckHygieneRejectsHelperBeforeTests(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

import "testing"

func helper(t *testing.T) {
	t.Helper()
}

func TestExample(t *testing.T) {}
`)

	violations := CheckHygiene(fileSet, file, "example_test.go")
	if !hasViolation(violations, analysis.DiagnosticTestHelperOrder) {
		t.Fatalf("expected test helper order violation, got %#v", violations)
	}
}

func TestCheckHygieneRejectsHelperBetweenTests(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

import "testing"

func TestFirst(t *testing.T) {}

func helper(t *testing.T) {
	t.Helper()
}

func TestSecond(t *testing.T) {}
`)

	violations := CheckHygiene(fileSet, file, "example_test.go")
	if !hasViolation(violations, analysis.DiagnosticTestHelperOrder) {
		t.Fatalf("expected test helper order violation, got %#v", violations)
	}
}

func TestCheckHygieneAllowsHelpersAfterTests(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

import "testing"

func TestExample(t *testing.T) {}

func helper(t *testing.T) {
	t.Helper()
}
`)

	violations := CheckHygiene(fileSet, file, "example_test.go")
	if hasViolation(violations, analysis.DiagnosticTestHelperOrder) {
		t.Fatalf("expected helper order to pass, got %#v", violations)
	}
}

func parseGoSource(t *testing.T, source string) (fileSet *token.FileSet, file *ast.File) {
	t.Helper()

	fileSet = token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "example.go", source, parser.ParseComments)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	return fileSet, file
}

func hasViolation(violations []analysis.Violation, rule string) (found bool) {
	for _, violation := range violations {
		if violation.Rule == rule {
			return true
		}
	}

	return false
}
