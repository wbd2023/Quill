package syntax

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
	"ciphera/tools/internal/checks/gopolicy"
)

func checkSecretErrorArguments(
	fileSet *token.FileSet,
	callExpression *ast.CallExpr,
	parameters gopolicy.ParameterConfig,
) (violations []analysis.Violation) {
	for _, argument := range callExpression.Args[1:] {
		if !expressionContainsSecretLikeIdentifier(argument, parameters.SecretNames) {
			continue
		}

		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(argument.Pos()),
			Rule:     analysis.DiagnosticErrorContextNoSecrets,
			Message:  "error context must not include secrets in fmt.Errorf arguments",
		})
	}

	return violations
}
