package checks

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/policy"
)

// CheckErrorHandlingStyle enforces Go error-message and sentinel-error style.
func CheckErrorHandlingStyle(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
	parameters policy.GoParameterConfig,
) (violations []Violation) {
	if !classifier.HasClass(path, PathClassGoSource) {
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

	if isTestFile || classifier.HasClass(path, PathClassDomainErrors) {
		return violations
	}

	return append(violations, collectSentinelErrorLocationViolations(fileSet, file, classifier)...)
}
