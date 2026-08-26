package syntax

import (
	"go/ast"
	"go/token"

	"github.com/wbd2023/quill/internal/checks/golang/analysis"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
)

/* --------------------------------------- Call Collection -------------------------------------- */

func collectErrorCallViolations(
	fileSet *token.FileSet,
	file *ast.File,
	isTestFile bool,
	parameters gopolicy.ParameterConfig,
	fmtImportAliases map[string]bool,
	errorsImportAliases map[string]bool,
) (violations []analysis.Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || len(call.Args) == 0 {
			return true
		}

		violations = append(violations, checkErrorCall(
			fileSet,
			call,
			isTestFile,
			parameters,
			fmtImportAliases,
			errorsImportAliases,
		)...)
		return true
	})

	return violations
}

/* ------------------------------------- Call Classification ------------------------------------ */

func checkErrorCall(
	fileSet *token.FileSet,
	call *ast.CallExpr,
	isTestFile bool,
	parameters gopolicy.ParameterConfig,
	fmtImportAliases map[string]bool,
	errorsImportAliases map[string]bool,
) (violations []analysis.Violation) {
	selectorExpression, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	ident, ok := selectorExpression.X.(*ast.Ident)
	if !ok {
		return nil
	}

	switch {
	case selectorExpression.Sel.Name == "Errorf" && fmtImportAliases[ident.Name]:
		return checkFmtErrorfCall(fileSet, call, isTestFile, parameters)
	case selectorExpression.Sel.Name == "New" && errorsImportAliases[ident.Name]:
		return checkErrorsNewCall(fileSet, call)
	}

	return nil
}

func checkFmtErrorfCall(
	fileSet *token.FileSet,
	call *ast.CallExpr,
	isTestFile bool,
	parameters gopolicy.ParameterConfig,
) (violations []analysis.Violation) {
	message, found := literalString(call.Args[0])
	if found {
		violations = append(
			violations,
			checkErrorMessageLiteralStyle(
				fileSet,
				call.Args[0],
				message,
				"fmt.Errorf",
			)...,
		)
	}

	if isTestFile {
		return violations
	}

	return append(violations, checkSecretErrorArguments(fileSet, call, parameters)...)
}

func checkErrorsNewCall(
	fileSet *token.FileSet,
	call *ast.CallExpr,
) (violations []analysis.Violation) {
	message, found := literalString(call.Args[0])
	if !found {
		return nil
	}

	return checkErrorMessageLiteralStyle(
		fileSet,
		call.Args[0],
		message,
		"errors.New",
	)
}
