package syntax

import (
	"go/ast"
	"go/token"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

func checkNamedStructLiterals(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []analysis.Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		literal, ok := node.(*ast.CompositeLit)
		if !ok || len(literal.Elts) == 0 {
			return true
		}

		if !isNamedStructLiteralType(literal.Type) {
			return true
		}

		for _, element := range literal.Elts {
			if _, ok := element.(*ast.KeyValueExpr); ok {
				return true
			}
		}

		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(literal.Pos()),
			Rule:     analysis.DiagnosticNamedStructLiterals,
			Message:  "struct literals must use named fields by default",
		})

		return true
	})

	return violations
}

func isNamedStructLiteralType(expression ast.Expr) (found bool) {
	switch expression.(type) {
	case *ast.Ident, *ast.SelectorExpr:
		return true
	default:
		return false
	}
}
