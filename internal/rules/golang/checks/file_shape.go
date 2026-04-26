package checks

import (
	"go/ast"
	"go/token"
)

const (
	maxHandwrittenFileLines = 300
	maxFunctionLines        = 80
	tinyGlueFileLines       = 15
)

// CheckFileShape reports objective file-granularity smells.
func CheckFileShape(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
) (violations []Violation) {
	lineCount := fileLineCount(fileSet, file)
	violations = append(violations, longFileViolations(fileSet, file, lineCount)...)
	violations = append(violations, vagueFileNameViolations(fileSet, file, path)...)
	violations = append(violations, tinyGlueFileViolations(fileSet, file, path, lineCount)...)

	if !isTestFile {
		violations = append(violations, longFunctionViolations(fileSet, file, path)...)
	}

	return violations
}
