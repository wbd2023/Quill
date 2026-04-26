package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

func longFileViolations(
	fileSet *token.FileSet,
	file *ast.File,
	lineCount int,
) (violations []Violation) {
	if lineCount <= maxHandwrittenFileLines {
		return nil
	}

	return []Violation{{
		Position: fileSet.Position(file.Package),
		Rule:     DiagnosticFileShapeLongFile,
		Message: fmt.Sprintf(
			"handwritten Go file has %d lines; split files over %d lines "+
				"unless the shape is clearly justified",
			lineCount,
			maxHandwrittenFileLines,
		),
	}}
}

func longFunctionViolations(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []Violation) {
	if longFunctionFileExempt(path) {
		return nil
	}

	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Body == nil {
			continue
		}

		lineCount := functionLineCount(fileSet, function)
		if lineCount <= maxFunctionLines {
			continue
		}

		violations = append(violations, Violation{
			Position: fileSet.Position(function.Pos()),
			Rule:     DiagnosticFileShapeLongFunction,
			Message: fmt.Sprintf(
				"function %s has %d lines; split functions over %d lines "+
					"unless they are parser or table-driven exceptions",
				function.Name.Name,
				lineCount,
				maxFunctionLines,
			),
		})
	}

	return violations
}

func fileLineCount(fileSet *token.FileSet, file *ast.File) (lines int) {
	start := fileSet.Position(file.Package).Line
	end := fileSet.Position(file.End()).Line
	if end < start {
		return 0
	}

	return end - start + 1
}

func longFunctionFileExempt(path string) (exempt bool) {
	name := strings.TrimSuffix(filepath.Base(path), ".go")
	return strings.Contains(name, "parse") || strings.Contains(name, "schema")
}

func functionLineCount(fileSet *token.FileSet, function *ast.FuncDecl) (lines int) {
	start := fileSet.Position(function.Pos()).Line
	end := fileSet.Position(function.End()).Line
	if end < start {
		return 0
	}

	return end - start + 1
}
