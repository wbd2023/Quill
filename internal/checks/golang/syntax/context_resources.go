package syntax

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
)

// CheckContextAndResourceSafety check context and resource safety.
func CheckContextAndResourceSafety(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier analysis.PathClassifier,
) (violations []analysis.Violation) {
	if !classifier.HasRole(path, analysis.PathRoleGoSource) {
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
