package structure

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/rules/golang/analysis"
)

func CheckGuardClauseSpacing(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []analysis.Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		block, ok := node.(*ast.BlockStmt)
		if !ok {
			return true
		}

		for index := 0; index+1 < len(block.List); index++ {
			firstIf, ok := block.List[index].(*ast.IfStmt)
			if !ok || !isGuardClause(firstIf) {
				continue
			}

			secondIf, ok := block.List[index+1].(*ast.IfStmt)
			if !ok || !isGuardClause(secondIf) {
				continue
			}

			firstCloseLine := fileSet.Position(firstIf.Body.Rbrace).Line
			secondStartLine := fileSet.Position(secondIf.Pos()).Line
			if secondStartLine-firstCloseLine != 1 {
				continue
			}

			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(secondIf.Pos()),
				Rule:     analysis.DiagnosticGuardClauseSpacing,
				Message:  "consecutive guard clauses should be separated by a blank line",
			})
		}

		return true
	})

	return violations
}

func isGuardClause(statement *ast.IfStmt) (guard bool) {
	if statement == nil ||
		statement.Else != nil ||
		statement.Body == nil ||
		len(statement.Body.List) == 0 {
		return false
	}

	_, ok := statement.Body.List[len(statement.Body.List)-1].(*ast.ReturnStmt)
	return ok
}
