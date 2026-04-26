package checks

import "go/ast"

func isBareErrReturn(returnStatement *ast.ReturnStmt) (found bool) {
	if len(returnStatement.Results) == 0 {
		return false
	}

	lastReturnExpression := returnStatement.Results[len(returnStatement.Results)-1]
	identifier, ok := lastReturnExpression.(*ast.Ident)
	if !ok {
		return false
	}

	return identifier.Name == "err"
}
