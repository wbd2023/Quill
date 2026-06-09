package structure

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"ciphera/tools/internal/checks/golang/analysis"
)

func parseGoSource(t *testing.T, source string) (fileSet *token.FileSet, file *ast.File) {
	t.Helper()

	fileSet = token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "example.go", source, parser.ParseComments)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	return fileSet, file
}

func sourceLines(source string) (lines []string) {
	return strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")
}

func hasViolation(violations []analysis.Violation, rule string) (found bool) {
	for _, violation := range violations {
		if violation.Rule == rule {
			return true
		}
	}

	return false
}

func hasViolationAt(
	violations []analysis.Violation,
	rule string,
	line int,
	messageFragment string,
) (found bool) {
	for _, violation := range violations {
		if violation.Rule != rule || violation.Position.Line != line {
			continue
		}

		if strings.Contains(violation.Message, messageFragment) {
			return true
		}
	}

	return false
}
