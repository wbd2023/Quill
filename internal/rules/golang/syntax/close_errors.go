package syntax

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/rules/golang/analysis"
)

func checkIgnoredCloseErrors(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []analysis.Violation) {
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

		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(assignStatement.Pos()),
			Rule:     analysis.DiagnosticIgnoredCloseErrorReason,
			Message:  "ignored close errors require an inline comment explaining why they are safe",
		})

		return true
	})

	return violations
}

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
