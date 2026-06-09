package syntax

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
)

func checkPointerToInterfaces(
	fileSet *token.FileSet,
	file *ast.File,
	interfaceNames map[string]bool,
) (violations []analysis.Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.Field:
			ruleID, message, found := pointerToInterfaceViolation(
				typed.Type,
				interfaceNames,
				true,
			)
			if !found {
				return true
			}

			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(typed.Type.Pos()),
				Rule:     ruleID,
				Message:  message,
			})

		case *ast.ValueSpec:
			ruleID, message, found := pointerToInterfaceViolation(
				typed.Type,
				interfaceNames,
				false,
			)
			if !found {
				return true
			}

			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(typed.Type.Pos()),
				Rule:     ruleID,
				Message:  message,
			})
		}

		return true
	})

	return violations
}

func pointerToInterfaceViolation(
	expression ast.Expr,
	interfaceNames map[string]bool,
	inField bool,
) (ruleID string, message string, found bool) {
	starExpression, ok := expression.(*ast.StarExpr)
	if !ok {
		return "", "", false
	}

	if !isInterfaceTypeExpression(starExpression.X, interfaceNames) {
		return "", "", false
	}

	if inField {
		return analysis.DiagnosticInterfaceValuesDirect,
			"functions and structs must pass interface values directly, not *interface types",
			true
	}

	return analysis.DiagnosticNoPointersToInterfaces,
		"pointers to interfaces are forbidden",
		true
}

func isInterfaceTypeExpression(expression ast.Expr, interfaceNames map[string]bool) (found bool) {
	switch typedExpression := expression.(type) {
	case *ast.InterfaceType:
		return true
	case *ast.Ident:
		return interfaceNames[typedExpression.Name]
	default:
		return false
	}
}
