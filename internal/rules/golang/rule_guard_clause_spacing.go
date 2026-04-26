package golang

import (
	"go/ast"
	"go/parser"
	"go/token"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rules/golang/checks"
)

func CheckGuardClauseSpacing(
	repoRoot string,
	directories []string,
	repository policy.RepositoryConfig,
) (result contract.ExecutionResult, err error) {
	files, err := goFilesInDirectories(directories, repository)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	fileSet := token.NewFileSet()
	diagnostics := make([]contract.Diagnostic, 0)

	for _, path := range files {
		syntax, parseErr := parser.ParseFile(fileSet, path, nil, parser.SkipObjectResolution)
		if parseErr != nil {
			return contract.ExecutionResult{}, parseErr
		}

		ast.Inspect(syntax, func(node ast.Node) bool {
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

				diagnostics = append(diagnostics, contract.Diagnostic{
					Code:    checks.DiagnosticGuardClauseSpacing,
					File:    relativePath(repoRoot, path),
					Line:    secondStartLine,
					Message: "consecutive guard clauses should be separated by a blank line",
				})
			}

			return true
		})
	}

	if len(diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return contract.ExecutionResult{
		Diagnostics: diagnostics,
	}, errViolationsFound
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
