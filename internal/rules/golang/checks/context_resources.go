package checks

import (
	"go/ast"
	"go/token"
)

func CheckContextAndResourceSafety(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
) (violations []Violation) {
	if !classifier.HasClass(path, PathClassGoSource) {
		return nil
	}

	contextAliases := importAliases(file, "context")
	httpAliases := importAliases(file, "net/http")

	if !isTestFile {
		violations = append(
			violations,
			checkContextFields(fileSet, file, contextAliases)...,
		)
		violations = append(
			violations,
			checkHTTPTimeouts(fileSet, file, httpAliases)...,
		)
	}

	violations = append(
		violations,
		checkIgnoredCloseErrors(fileSet, file)...,
	)

	return violations
}
