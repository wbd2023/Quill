package structure

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
)

const (
	maxHandwrittenFileLines = 300
	maxFunctionLines        = 80
	tinyGlueFileLines       = 15
)

// CheckShape reports objective file-granularity smells.
func CheckShape(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
) (violations []analysis.Violation) {
	lineCount := fileLineCount(fileSet, file)
	violations = append(violations, longFileViolations(fileSet, file, lineCount)...)
	violations = append(violations, vagueFileNameViolations(fileSet, file, path)...)
	violations = append(violations, tinyGlueFileViolations(fileSet, file, path, lineCount)...)

	if !isTestFile {
		violations = append(violations, longFunctionViolations(fileSet, file, path)...)
	}

	return violations
}
