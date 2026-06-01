package structure

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"ciphera/tools/internal/rules/golang/analysis"
)

func CheckScannerEntrypointOrder(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []analysis.Violation) {
	if !isFirstPartyScannerFile(path) {
		return nil
	}

	pendingHelpers := make([]*ast.FuncDecl, 0)
	foundEntrypoint := false
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Recv != nil {
			continue
		}

		if isScannerEntrypoint(function.Name.Name) {
			foundEntrypoint = true
			for _, helper := range pendingHelpers {
				violations = append(violations, scannerEntrypointViolation(fileSet, helper))
			}
			pendingHelpers = nil
			continue
		}

		if !foundEntrypoint {
			pendingHelpers = append(pendingHelpers, function)
		}
	}

	if !foundEntrypoint {
		return nil
	}

	return violations
}

func scannerEntrypointViolation(
	fileSet *token.FileSet,
	function *ast.FuncDecl,
) (violation analysis.Violation) {
	return analysis.Violation{
		Position: fileSet.Position(function.Pos()),
		Rule:     analysis.DiagnosticScannerEntrypointOrder,
		Message: fmt.Sprintf(
			"scanner helper %q appears before exported Check... entrypoint",
			function.Name.Name,
		),
	}
}

func isFirstPartyScannerFile(path string) (found bool) {
	if strings.HasSuffix(path, "_test.go") {
		return false
	}

	for _, directory := range []string{
		"/tools/internal/rules/bash/",
		"/tools/internal/rules/vocabulary/",
		"/tools/internal/rules/security/",
		"/tools/internal/rules/text/",
	} {
		if strings.Contains(path, directory) {
			return true
		}
	}

	return false
}

func isScannerEntrypoint(name string) (found bool) {
	return strings.HasPrefix(name, "Check") && ast.IsExported(name)
}
