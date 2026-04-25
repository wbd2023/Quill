package gostyle

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	repostyle "ciphera/tools/internal/rules/repo"
)

func CheckGuardClauseSpacing(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	files, err := repostyle.CollectFiles(repoRoot, repository, scope, ".go")
	if err != nil {
		return "", err
	}

	fileSet := token.NewFileSet()
	var builder strings.Builder
	found := false

	for _, path := range files {
		parsedFile, parseErr := parser.ParseFile(fileSet, path, nil, parser.SkipObjectResolution)
		if parseErr != nil {
			return "", parseErr
		}

		ast.Inspect(parsedFile, func(node ast.Node) bool {
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

				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d consecutive guard clauses should be separated by a blank line\n",
					repostyle.RelativePath(repoRoot, path),
					secondStartLine,
				))
			}

			return true
		})
	}

	if !found {
		return "", nil
	}

	return builder.String(), errViolationsFound
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
