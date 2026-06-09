package syntax

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
)

/* ---------------------------------------- HTTP Timeouts --------------------------------------- */

func checkHTTPTimeouts(
	fileSet *token.FileSet,
	file *ast.File,
	httpAliases map[string]bool,
) (violations []analysis.Violation) {
	if len(httpAliases) == 0 {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		switch typedNode := node.(type) {
		case *ast.CallExpr:
			selector, ok := typedNode.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			packageIdentifier, ok := selector.X.(*ast.Ident)
			if !ok || !httpAliases[packageIdentifier.Name] {
				return true
			}

			switch selector.Sel.Name {
			case "Get", "Post", "PostForm", "Head":
				violations = append(violations, analysis.Violation{
					Position: fileSet.Position(typedNode.Pos()),
					Rule:     analysis.DiagnosticExplicitNetworkTimeouts,
					Message:  "network requests must use an http.Client with an explicit timeout",
				})
			}

		case *ast.SelectorExpr:
			packageIdentifier, ok := typedNode.X.(*ast.Ident)
			if !ok || !httpAliases[packageIdentifier.Name] {
				return true
			}
			if typedNode.Sel.Name != "DefaultClient" {
				return true
			}

			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(typedNode.Pos()),
				Rule:     analysis.DiagnosticExplicitNetworkTimeouts,
				Message:  "network clients must not rely on http.DefaultClient",
			})

		case *ast.CompositeLit:
			if !isHTTPClientType(typedNode.Type, httpAliases) ||
				httpClientHasTimeoutField(typedNode) {
				return true
			}

			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(typedNode.Pos()),
				Rule:     analysis.DiagnosticExplicitNetworkTimeouts,
				Message:  "http.Client literals must set Timeout explicitly",
			})
		}

		return true
	})

	return violations
}

/* ----------------------------------------- Predicates ----------------------------------------- */

func isHTTPClientType(expression ast.Expr, httpAliases map[string]bool) (found bool) {
	selector, ok := expression.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Client" {
		return false
	}

	packageIdentifier, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	return httpAliases[packageIdentifier.Name]
}

func httpClientHasTimeoutField(literal *ast.CompositeLit) (found bool) {
	for _, element := range literal.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		keyIdentifier, ok := field.Key.(*ast.Ident)
		if ok && keyIdentifier.Name == "Timeout" {
			return true
		}
	}

	return false
}
