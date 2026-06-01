package syntax

import (
	"go/ast"
)

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
