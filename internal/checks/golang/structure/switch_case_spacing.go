package structure

import (
	"go/ast"
	"go/token"
	"strings"

	"ciphera/tools/internal/checks/golang/analysis"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// switch_case_spacing constants.
const (
	minSwitchCaseCount        = 2
	maxVerySmallSwitchClauses = 3
	maxVerySmallCaseSpan      = 4
)

/* ---------------------------------------- Spacing Rules --------------------------------------- */

// CheckSwitchCaseSpacing checks that non-trivial switch statements separate case blocks with blank
// lines. Blank lines between cases are always acceptable; they are never flagged.
func CheckSwitchCaseSpacing(
	fileSet *token.FileSet,
	file *ast.File,
	lines []string,
) (violations []analysis.Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch typedNode := node.(type) {
		case *ast.SwitchStmt:
			violations = append(violations, switchSpacingViolations(
				fileSet,
				lines,
				typedNode.Body.List,
			)...)

		case *ast.TypeSwitchStmt:
			violations = append(violations, switchSpacingViolations(
				fileSet,
				lines,
				typedNode.Body.List,
			)...)
		}

		return true
	})

	return violations
}

/* --------------------------------------- Switch Analysis -------------------------------------- */

func switchSpacingViolations(
	fileSet *token.FileSet,
	lines []string,
	caseStatements []ast.Stmt,
) (violations []analysis.Violation) {
	caseClauses := collectCaseClauses(caseStatements)
	if len(caseClauses) < minSwitchCaseCount {
		return nil
	}

	if !isNonTrivialSwitch(fileSet, caseClauses) {
		return nil
	}

	for index := 0; index+1 < len(caseClauses); index++ {
		previousClause := caseClauses[index]
		nextClause := caseClauses[index+1]

		if hasBlankLineBetween(fileSet, lines, previousClause, nextClause) {
			continue
		}

		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(nextClause.Pos()),
			Rule:     analysis.DiagnosticSwitchCaseSpacing,
			Message: "non-trivial switch statements should separate case blocks " +
				"with a blank line",
		})
	}

	return violations
}

/* ------------------------------------- Case Classification ------------------------------------ */

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
