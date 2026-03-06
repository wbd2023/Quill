package support

import (
	"go/ast"
	"go/token"
	"strconv"
)

/* -------------------------------------- AST Type Helpers -------------------------------------- */

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

// TypeString extracts a human-readable type name from an AST expression.
func TypeString(expression ast.Expr) (typeName string) {
	switch typed := expression.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.SelectorExpr:
		return TypeString(typed.X) + "." + typed.Sel.Name
	case *ast.StarExpr:
		return TypeString(typed.X)
	case *ast.ArrayType:
		return "[]" + TypeString(typed.Elt)
	default:
		return ""
	}
}

func ExtractStringLiteral(expression ast.Expr) (value string, found bool) {
	literal, ok := expression.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return "", false
	}

	unquotedValue, err := strconv.Unquote(literal.Value)
	if err != nil {
		return "", false
	}

	return unquotedValue, true
}
