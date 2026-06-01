package structure

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"ciphera/tools/internal/rules/golang/analysis"
)

func tinyGlueFileViolations(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	lineCount int,
) (violations []analysis.Violation) {
	if !tinyGlueFileNeedsMerge(file, path, lineCount) {
		return nil
	}

	return []analysis.Violation{{
		Position: fileSet.Position(file.Package),
		Rule:     analysis.DiagnosticFileShapeTinyGlue,
		Message: fmt.Sprintf(
			"tiny file has %d lines and only one unexported helper; "+
				"merge it unless it removes real navigation cost",
			lineCount,
		),
	}}
}

func tinyGlueFileNeedsMerge(file *ast.File, path string, lineCount int) (needsMerge bool) {
	if lineCount >= tinyGlueFileLines || isTinyFileExempt(path) || len(file.Decls) != 1 {
		return false
	}

	switch declaration := file.Decls[0].(type) {
	case *ast.FuncDecl:
		return unexportedHelperFunction(declaration)
	case *ast.GenDecl:
		return unexportedTypeAlias(declaration)
	}

	return false
}

func isTinyFileExempt(path string) (exempt bool) {
	name := filepath.Base(path)
	return name == "doc.go" || name == "main.go"
}

func unexportedHelperFunction(function *ast.FuncDecl) (helper bool) {
	name := function.Name.Name
	if name == "init" || name == "main" {
		return false
	}

	return !ast.IsExported(name)
}

func unexportedTypeAlias(declaration *ast.GenDecl) (alias bool) {
	if declaration.Tok != token.TYPE || len(declaration.Specs) != 1 {
		return false
	}

	spec, ok := declaration.Specs[0].(*ast.TypeSpec)
	return ok && spec.Assign.IsValid() && !ast.IsExported(spec.Name.Name)
}
