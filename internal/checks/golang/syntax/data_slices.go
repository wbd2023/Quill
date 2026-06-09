package syntax

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
)

/* ---------------------------------------- Slice Checks ---------------------------------------- */

func checkSliceEmptinessStyle(
	fileSet *token.FileSet,
	file *ast.File,
	sliceNames map[string]bool,
) (violations []analysis.Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		binaryExpression, ok := node.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		switch binaryExpression.Op {
		case token.LOR, token.LAND:
		default:
			return true
		}

		name, found := sliceNilGuardPair(binaryExpression.X, binaryExpression.Y, sliceNames)
		if !found {
			name, found = sliceNilGuardPair(binaryExpression.Y, binaryExpression.X, sliceNames)
		}
		if !found {
			return true
		}

		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(binaryExpression.Pos()),
			Rule:     analysis.DiagnosticLenForSliceEmptiness,
			Message:  "slice emptiness checks must use len(" + name + ") instead of nil guards",
		})

		return true
	})

	return violations
}

/* ----------------------------------------- Nil Guards ----------------------------------------- */

func sliceNilGuardPair(
	nilSide ast.Expr,
	lenSide ast.Expr,
	sliceNames map[string]bool,
) (name string, found bool) {
	name, nilCompared, found := nilComparedSliceName(nilSide, sliceNames)
	if !found || !nilCompared {
		return "", false
	}

	if !containsSliceLenCheck(lenSide, name) {
		return "", false
	}

	return name, true
}

func nilComparedSliceName(
	expression ast.Expr,
	sliceNames map[string]bool,
) (name string, nilCompared bool, found bool) {
	binaryExpression, ok := expression.(*ast.BinaryExpr)
	if !ok {
		return "", false, false
	}

	switch binaryExpression.Op {
	case token.EQL, token.NEQ:
	default:
		return "", false, false
	}

	if identifier, ok := binaryExpression.X.(*ast.Ident); ok &&
		isNilIdentifier(binaryExpression.Y) {
		if sliceNames[identifier.Name] {
			return identifier.Name, true, true
		}
	}
	if identifier, ok := binaryExpression.Y.(*ast.Ident); ok &&
		isNilIdentifier(binaryExpression.X) {
		if sliceNames[identifier.Name] {
			return identifier.Name, true, true
		}
	}

	return "", false, false
}

func containsSliceLenCheck(expression ast.Expr, name string) (found bool) {
	binaryExpression, ok := expression.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	callExpression, ok := binaryExpression.X.(*ast.CallExpr)
	if !ok || !isLenCallForName(callExpression, name) {
		callExpression, ok = binaryExpression.Y.(*ast.CallExpr)
		if !ok || !isLenCallForName(callExpression, name) {
			return false
		}
	}

	return true
}

func isLenCallForName(callExpression *ast.CallExpr, name string) (found bool) {
	functionName, ok := callExpression.Fun.(*ast.Ident)
	if !ok || functionName.Name != "len" || len(callExpression.Args) != 1 {
		return false
	}

	identifier, ok := callExpression.Args[0].(*ast.Ident)
	return ok && identifier.Name == name
}

func isNilIdentifier(expression ast.Expr) (found bool) {
	identifier, ok := expression.(*ast.Ident)
	return ok && identifier.Name == "nil"
}
