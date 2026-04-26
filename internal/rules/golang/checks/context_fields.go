package checks

import (
	"go/ast"
	"go/token"
)

func checkContextFields(
	fileSet *token.FileSet,
	file *ast.File,
	contextAliases map[string]bool,
) (violations []Violation) {
	if len(contextAliases) == 0 {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		typeSpec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok || structType.Fields == nil {
			return true
		}

		for _, field := range structType.Fields.List {
			if !isContextType(field.Type, contextAliases) {
				continue
			}

			violations = append(violations, Violation{
				Position: fileSet.Position(field.Pos()),
				Rule:     DiagnosticNoContextOnStructs,
				Message:  "contexts must not be stored on struct fields",
			})
		}

		return true
	})

	return violations
}

func isContextType(expression ast.Expr, contextAliases map[string]bool) (found bool) {
	selector, ok := expression.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Context" {
		return false
	}

	packageIdentifier, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	return contextAliases[packageIdentifier.Name]
}
