package checks

import (
	"go/ast"
	"go/token"
)

func CheckDataUsage(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
) (violations []Violation) {
	if !classifier.HasClass(path, PathClassGoSource) {
		return nil
	}

	interfaceNames := collectInterfaceTypeNames(file)
	sliceNames := collectSliceNames(file)

	if !isTestFile {
		violations = append(
			violations,
			checkNamedStructLiterals(fileSet, file)...,
		)
	}

	violations = append(
		violations,
		checkPointerToInterfaces(fileSet, file, interfaceNames)...,
	)
	violations = append(
		violations,
		checkSliceEmptinessStyle(fileSet, file, sliceNames)...,
	)

	return violations
}
