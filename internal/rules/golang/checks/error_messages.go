package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func checkErrorMessageLiteralStyle(
	fileSet *token.FileSet,
	expression ast.Expr,
	message string,
	callName string,
) (violations []Violation) {
	trimmedMessage := strings.TrimSpace(message)
	if trimmedMessage == "" {
		return nil
	}

	if startsWithUppercaseLetter(trimmedMessage) {
		violations = append(violations, Violation{
			Position: fileSet.Position(expression.Pos()),
			Rule:     DiagnosticErrorContextLowercase,
			Message:  fmt.Sprintf("error context must be lowercase (%s)", callName),
		})
	}

	if endsWithSentencePunctuation(trimmedMessage) {
		violations = append(violations, Violation{
			Position: fileSet.Position(expression.Pos()),
			Rule:     DiagnosticErrorContextNoPunctuation,
			Message:  fmt.Sprintf("error context must not end with punctuation (%s)", callName),
		})
	}

	return violations
}
