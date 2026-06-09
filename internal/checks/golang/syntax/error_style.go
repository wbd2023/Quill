package syntax

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
	gopolicy "ciphera/tools/internal/checks/golang/policy"
)

// CheckErrorHandlingStyle enforces Go error-message and sentinel-error style.
func CheckErrorHandlingStyle(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier analysis.PathClassifier,
	parameters gopolicy.ParameterConfig,
) (violations []analysis.Violation) {
	if !classifier.HasRole(path, analysis.PathRoleGoSource) {
		return nil
	}

	fmtImportAliases := importAliases(file, "fmt")
	errorsImportAliases := importAliases(file, "errors")
	violations = append(violations, collectErrorCallViolations(
		fileSet,
		file,
		isTestFile,
		parameters,
		fmtImportAliases,
		errorsImportAliases,
	)...)

	if isTestFile || classifier.HasRole(path, analysis.PathRoleDomainErrors) {
		return violations
	}

	return append(violations, collectSentinelErrorLocationViolations(fileSet, file, classifier)...)
}
