package syntax

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"ciphera/tools/internal/rules/golang/analysis"
)

func checkErrorMessageLiteralStyle(
	fileSet *token.FileSet,
	expression ast.Expr,
	message string,
	callName string,
) (violations []analysis.Violation) {
	trimmedMessage := strings.TrimSpace(message)
	if trimmedMessage == "" {
		return nil
	}

	if startsWithUppercaseLetter(trimmedMessage) {
		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(expression.Pos()),
			Rule:     analysis.DiagnosticErrorContextLowercase,
			Message:  fmt.Sprintf("error context must be lowercase (%s)", callName),
		})
	}

	if endsWithSentencePunctuation(trimmedMessage) {
		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(expression.Pos()),
			Rule:     analysis.DiagnosticErrorContextNoPunctuation,
			Message:  fmt.Sprintf("error context must not end with punctuation (%s)", callName),
		})
	}

	return violations
}
