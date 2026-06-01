package syntax

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/rules/golang/analysis"
)

func CheckDataUsage(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier analysis.PathClassifier,
) (violations []analysis.Violation) {
	if !classifier.HasRole(path, analysis.PathRoleGoSource) {
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
