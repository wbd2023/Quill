package structure

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
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
