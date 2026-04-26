package checks

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/policy"
)

/* --------------------------------------- Call Collection -------------------------------------- */

func collectErrorCallViolations(
	fileSet *token.FileSet,
	file *ast.File,
	isTestFile bool,
	parameters policy.GoParameterConfig,
	fmtImportAliases map[string]bool,
	errorsImportAliases map[string]bool,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok || len(callExpression.Args) == 0 {
			return true
		}

		violations = append(violations, checkErrorCall(
			fileSet,
			callExpression,
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
	callExpression *ast.CallExpr,
	isTestFile bool,
	parameters policy.GoParameterConfig,
	fmtImportAliases map[string]bool,
	errorsImportAliases map[string]bool,
) (violations []Violation) {
	selectorExpression, ok := callExpression.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	packageIdentifier, ok := selectorExpression.X.(*ast.Ident)
	if !ok {
		return nil
	}

	switch {
	case selectorExpression.Sel.Name == "Errorf" && fmtImportAliases[packageIdentifier.Name]:
		return checkFmtErrorfCall(fileSet, callExpression, isTestFile, parameters)
	case selectorExpression.Sel.Name == "New" && errorsImportAliases[packageIdentifier.Name]:
		return checkErrorsNewCall(fileSet, callExpression)
	}

	return nil
}

func checkFmtErrorfCall(
	fileSet *token.FileSet,
	callExpression *ast.CallExpr,
	isTestFile bool,
	parameters policy.GoParameterConfig,
) (violations []Violation) {
	message, found := literalString(callExpression.Args[0])
	if found {
		violations = append(
			violations,
			checkErrorMessageLiteralStyle(
				fileSet,
				callExpression.Args[0],
				message,
				"fmt.Errorf",
			)...,
		)
	}

	if isTestFile {
		return violations
	}

	return append(violations, checkSecretErrorArguments(fileSet, callExpression, parameters)...)
}

func checkErrorsNewCall(
	fileSet *token.FileSet,
	callExpression *ast.CallExpr,
) (violations []Violation) {
	message, found := literalString(callExpression.Args[0])
	if !found {
		return nil
	}

	return checkErrorMessageLiteralStyle(
		fileSet,
		callExpression.Args[0],
		message,
		"errors.New",
	)
}
