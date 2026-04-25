package gostyle

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	repostyle "ciphera/tools/internal/rules/repo"
)

const (
	minSwitchCaseCount        = 2
	maxVerySmallSwitchClauses = 3
	maxVerySmallCaseSpan      = 4
)

/* --------------------------------------- Switch Spacing --------------------------------------- */

func CheckSwitchCaseSpacing(
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
		contents, readErr := os.ReadFile(path)
		if readErr != nil {
			return "", readErr
		}

		lines := strings.Split(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\n")
		parsedFile, parseErr := parser.ParseFile(
			fileSet,
			path,
			contents,
			parser.SkipObjectResolution,
		)
		if parseErr != nil {
			return "", parseErr
		}

		ast.Inspect(parsedFile, func(node ast.Node) bool {
			switch typedNode := node.(type) {
			case *ast.SwitchStmt:
				recordSwitchSpacingViolations(
					repoRoot,
					path,
					fileSet,
					lines,
					typedNode.Body.List,
					&builder,
					&found,
				)

			case *ast.TypeSwitchStmt:
				recordSwitchSpacingViolations(
					repoRoot,
					path,
					fileSet,
					lines,
					typedNode.Body.List,
					&builder,
					&found,
				)
			}

			return true
		})
	}

	if !found {
		return "", nil
	}

	return builder.String(), errViolationsFound
}

func recordSwitchSpacingViolations(
	repoRoot string,
	path string,
	fileSet *token.FileSet,
	lines []string,
	caseStatements []ast.Stmt,
	builder *strings.Builder,
	found *bool,
) {
	caseClauses := collectCaseClauses(caseStatements)
	if len(caseClauses) < minSwitchCaseCount {
		return
	}

	verySmall := isVerySmallSwitch(fileSet, caseClauses)
	nonTrivial := isNonTrivialSwitch(fileSet, caseClauses)
	if !verySmall && !nonTrivial {
		return
	}

	for index := 0; index+1 < len(caseClauses); index++ {
		previousClause := caseClauses[index]
		nextClause := caseClauses[index+1]
		nextLine := fileSet.Position(nextClause.Pos()).Line

		if hasBlankLineBetween(fileSet, lines, previousClause, nextClause) {
			if verySmall {
				*found = true
				_, _ = fmt.Fprintf(
					builder,
					"%s:%d very small switch statements should stay compact "+
						"without blank lines between case blocks\n",
					repostyle.RelativePath(repoRoot, path),
					nextLine,
				)
			}

			continue
		}

		if nonTrivial {
			*found = true
			_, _ = fmt.Fprintf(
				builder,
				"%s:%d non-trivial switch statements should separate case blocks "+
					"with a blank line\n",
				repostyle.RelativePath(repoRoot, path),
				nextLine,
			)
		}
	}
}

func collectCaseClauses(statements []ast.Stmt) (clauses []*ast.CaseClause) {
	clauses = make([]*ast.CaseClause, 0, len(statements))
	for _, statement := range statements {
		caseClause, ok := statement.(*ast.CaseClause)
		if !ok {
			continue
		}

		clauses = append(clauses, caseClause)
	}

	return clauses
}

func isVerySmallSwitch(fileSet *token.FileSet, clauses []*ast.CaseClause) (found bool) {
	if len(clauses) < minSwitchCaseCount || len(clauses) > maxVerySmallSwitchClauses {
		return false
	}

	for _, clause := range clauses {
		if !isVerySmallCaseClause(fileSet, clause) {
			return false
		}
	}

	return true
}

func isNonTrivialSwitch(fileSet *token.FileSet, clauses []*ast.CaseClause) (found bool) {
	if len(clauses) > maxVerySmallSwitchClauses {
		return true
	}

	for _, clause := range clauses {
		if !isVerySmallCaseClause(fileSet, clause) {
			return true
		}
	}

	return false
}

func isVerySmallCaseClause(fileSet *token.FileSet, clause *ast.CaseClause) (found bool) {
	if clause == nil || len(clause.Body) > 1 {
		return false
	}

	endLine := fileSet.Position(clause.Colon).Line
	if len(clause.Body) == 1 {
		endLine = fileSet.Position(clause.Body[0].End()).Line
	}

	startLine := fileSet.Position(clause.Pos()).Line
	return endLine-startLine <= maxVerySmallCaseSpan
}

func hasBlankLineBetween(
	fileSet *token.FileSet,
	lines []string,
	previousClause *ast.CaseClause,
	nextClause *ast.CaseClause,
) (found bool) {
	previousEndLine := fileSet.Position(previousClause.End()).Line
	nextStartLine := fileSet.Position(nextClause.Pos()).Line
	if nextStartLine <= previousEndLine+1 {
		return false
	}

	for lineIndex := previousEndLine; lineIndex < nextStartLine-1; lineIndex++ {
		if lineIndex < 0 || lineIndex >= len(lines) {
			continue
		}

		if strings.TrimSpace(lines[lineIndex]) == "" {
			return true
		}
	}

	return false
}
