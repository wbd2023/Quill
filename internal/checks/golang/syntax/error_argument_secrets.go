package syntax

import (
	"go/ast"
	"go/token"

	"github.com/wbd2023/quill/internal/checks/golang/analysis"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
)

func checkSecretErrorArguments(
	fileSet *token.FileSet,
	call *ast.CallExpr,
	parameters gopolicy.ParameterConfig,
) (violations []analysis.Violation) {
	for _, argument := range call.Args[1:] {
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
