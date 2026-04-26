package checks

import (
	"go/ast"
	"go/token"

	"ciphera/tools/internal/policy"
)

func checkSecretErrorArguments(
	fileSet *token.FileSet,
	callExpression *ast.CallExpr,
	parameters policy.GoParameterConfig,
) (violations []Violation) {
	for _, argument := range callExpression.Args[1:] {
		if !expressionContainsSecretLikeIdentifier(argument, parameters.SecretNames) {
			continue
		}

		violations = append(violations, Violation{
			Position: fileSet.Position(argument.Pos()),
			Rule:     DiagnosticErrorContextNoSecrets,
			Message:  "error context must not include secrets in fmt.Errorf arguments",
		})
	}

	return violations
}
