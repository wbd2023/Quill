package checks

import (
	"go/ast"
	"go/token"
)

/* ------------------------------------- AST Type Extraction ------------------------------------ */

// TypeNameFromExpr extracts the named type from a basic AST expression.
func TypeNameFromExpr(expression ast.Expr) (name string) {
	switch typed := expression.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.SelectorExpr:
		return typed.Sel.Name
	default:
		return ""
	}
}

// ImplementationTypeFromAssertion resolves concrete implementation names from assertions.
func ImplementationTypeFromAssertion(expression ast.Expr) (name string) {
	switch typed := expression.(type) {
	case *ast.CallExpr:
		return ImplementationTypeFromAssertion(typed.Fun)

	case *ast.ParenExpr:
		return ImplementationTypeFromAssertion(typed.X)

	case *ast.StarExpr:
		return TypeNameFromExpr(typed.X)

	case *ast.UnaryExpr:
		if typed.Op == token.AND {
			return ImplementationTypeFromAssertion(typed.X)
		}
		return ""

	case *ast.CompositeLit:
		return TypeNameFromExpr(typed.Type)

	case *ast.Ident:
		return typed.Name

	default:
		return ""
	}
}

// ReceiverTypeName returns the receiver type for methods (supports T and *T).
func ReceiverTypeName(expression ast.Expr) (typeName string) {
	switch typed := expression.(type) {
	case *ast.StarExpr:
		if identifierNode, ok := typed.X.(*ast.Ident); ok {
			return identifierNode.Name
		}
	case *ast.Ident:
		return typed.Name
	}
	return ""
}

/* --------------------------------------- Type Formatting -------------------------------------- */

// typeString extracts a human-readable type name from an AST expression.
func typeString(expression ast.Expr) (typeName string) {
	switch typed := expression.(type) {
	case *ast.Ident:
		return typed.Name

	case *ast.SelectorExpr:
		return typeString(typed.X) + "." + typed.Sel.Name

	case *ast.StarExpr:
		return typeString(typed.X)

	case *ast.ArrayType:
		return "[]" + typeString(typed.Elt)

	default:
		return ""
	}
}
