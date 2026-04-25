package checks

import (
	"go/ast"
	"go/token"
	"path"
	"strconv"
)

func importAliases(file *ast.File, importPath string) (aliases map[string]bool) {
	aliases = make(map[string]bool)
	defaultAlias := path.Base(importPath)

	for _, importSpec := range file.Imports {
		pathValue, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil || pathValue != importPath {
			continue
		}

		if importSpec.Name != nil {
			if importSpec.Name.Name == "_" || importSpec.Name.Name == "." {
				continue
			}

			aliases[importSpec.Name.Name] = true
			continue
		}

		aliases[defaultAlias] = true
	}

	return aliases
}

func stringLiteralValue(expression ast.Expr) (value string, found bool) {
	literal, ok := expression.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return "", false
	}

	value, err := strconv.Unquote(literal.Value)
	if err != nil {
		return "", false
	}

	return value, true
}

func rightmostName(expression ast.Expr) (name string, found bool) {
	switch typedExpression := expression.(type) {
	case *ast.Ident:
		return typedExpression.Name, true
	case *ast.SelectorExpr:
		return typedExpression.Sel.Name, true
	default:
		return "", false
	}
}
