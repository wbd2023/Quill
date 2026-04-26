package checks

import (
	"go/ast"
	"go/token"
)

/* --------------------------------- Context and Resource Rules --------------------------------- */

func CheckContextAndResourceSafety(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
) (violations []Violation) {
	if !classifier.HasClass(path, PathClassGoSource) {
		return nil
	}

	contextAliases := importAliases(file, "context")
	httpAliases := importAliases(file, "net/http")

	if !isTestFile {
		violations = append(
			violations,
			checkContextFields(fileSet, file, contextAliases)...,
		)
		violations = append(
			violations,
			checkHTTPTimeouts(fileSet, file, httpAliases)...,
		)
	}

	violations = append(
		violations,
		checkIgnoredCloseErrors(fileSet, file)...,
	)

	return violations
}

/* --------------------------------------- Context Fields --------------------------------------- */

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

/* ---------------------------------------- HTTP Timeouts --------------------------------------- */

func checkHTTPTimeouts(
	fileSet *token.FileSet,
	file *ast.File,
	httpAliases map[string]bool,
) (violations []Violation) {
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
				violations = append(violations, Violation{
					Position: fileSet.Position(typedNode.Pos()),
					Rule:     DiagnosticExplicitNetworkTimeouts,
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

			violations = append(violations, Violation{
				Position: fileSet.Position(typedNode.Pos()),
				Rule:     DiagnosticExplicitNetworkTimeouts,
				Message:  "network clients must not rely on http.DefaultClient",
			})

		case *ast.CompositeLit:
			if !isHTTPClientType(typedNode.Type, httpAliases) ||
				httpClientHasTimeoutField(typedNode) {
				return true
			}

			violations = append(violations, Violation{
				Position: fileSet.Position(typedNode.Pos()),
				Rule:     DiagnosticExplicitNetworkTimeouts,
				Message:  "http.Client literals must set Timeout explicitly",
			})
		}

		return true
	})

	return violations
}

/* ------------------------------------ Ignored Close Errors ------------------------------------ */

func checkIgnoredCloseErrors(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		assignStatement, ok := node.(*ast.AssignStmt)
		if !ok || len(assignStatement.Lhs) != 1 || len(assignStatement.Rhs) != 1 {
			return true
		}

		blankIdentifier, ok := assignStatement.Lhs[0].(*ast.Ident)
		if !ok || blankIdentifier.Name != "_" || !isCloseCall(assignStatement.Rhs[0]) {
			return true
		}

		if hasCloseIgnoreJustification(fileSet, file, assignStatement.Pos()) {
			return true
		}

		violations = append(violations, Violation{
			Position: fileSet.Position(assignStatement.Pos()),
			Rule:     DiagnosticIgnoredCloseErrorReason,
			Message:  "ignored close errors require an inline comment explaining why they are safe",
		})

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

/* ------------------------------------ Close Call Detection ------------------------------------ */

func isCloseCall(expression ast.Expr) (found bool) {
	callExpression, ok := expression.(*ast.CallExpr)
	if !ok {
		return false
	}

	selector, ok := callExpression.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	return selector.Sel.Name == "Close"
}

func hasCloseIgnoreJustification(
	fileSet *token.FileSet,
	file *ast.File,
	position token.Pos,
) (found bool) {
	statementLine := fileSet.Position(position).Line

	for _, commentGroup := range file.Comments {
		startLine := fileSet.Position(commentGroup.Pos()).Line
		endLine := fileSet.Position(commentGroup.End()).Line
		if startLine > statementLine || endLine < statementLine-1 {
			continue
		}

		return true
	}

	return false
}
