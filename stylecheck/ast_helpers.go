package main

import (
	"go/ast"
	"go/token"
)

/* -------------------------------------- AST Type Helpers -------------------------------------- */

func typeNameFromExpr(expression ast.Expr) (name string) {
	switch typed := expression.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.SelectorExpr:
		return typed.Sel.Name
	default:
		return ""
	}
}

func implementationTypeFromAssertion(expression ast.Expr) (name string) {
	switch typed := expression.(type) {
	case *ast.CallExpr:
		return implementationTypeFromAssertion(typed.Fun)

	case *ast.ParenExpr:
		return implementationTypeFromAssertion(typed.X)

	case *ast.StarExpr:
		return typeNameFromExpr(typed.X)

	case *ast.UnaryExpr:
		if typed.Op == token.AND {
			return implementationTypeFromAssertion(typed.X)
		}
		return ""

	case *ast.CompositeLit:
		return typeNameFromExpr(typed.Type)

	case *ast.Ident:
		return typed.Name

	default:
		return ""
	}
}

// receiverTypeName returns the receiver type for methods (supports T and *T).
func receiverTypeName(expression ast.Expr) (typeName string) {
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
