package checks

import (
	"fmt"
	"go/ast"
	"go/token"

	"ciphera/tools/internal/policy"
)

/* ------------------------------------- Structured Logging ------------------------------------- */

func CheckStructuredLogging(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	classifier PathClassifier,
	parameters policy.GoParameterConfig,
) (violations []Violation) {
	if !classifier.HasClass(path, PathClassGoSource) {
		return nil
	}

	slogAliases := importAliases(file, "log/slog")
	if len(slogAliases) == 0 {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		violations = append(violations, checkStructuredLogCall(
			fileSet,
			callExpression,
			slogAliases,
			parameters,
		)...)

		return true
	})

	return violations
}

func checkStructuredLogCall(
	fileSet *token.FileSet,
	callExpression *ast.CallExpr,
	slogAliases map[string]bool,
	parameters policy.GoParameterConfig,
) (violations []Violation) {
	selector, ok := callExpression.Fun.(*ast.SelectorExpr)
	if !ok || !isStructuredLogCall(selector, slogAliases) {
		return nil
	}

	fieldCount := len(callExpression.Args) - 1
	if fieldCount <= 0 {
		return nil
	}

	if fieldCount%2 != 0 {
		return []Violation{{
			Position: fileSet.Position(callExpression.Pos()),
			Rule:     DiagnosticStructuredLogs,
			Message:  "structured log calls must use key/value pairs",
		}}
	}

	for argumentIndex := 1; argumentIndex < len(callExpression.Args); argumentIndex += 2 {
		violations = append(violations, checkStructuredLogPair(
			fileSet,
			callExpression.Args[argumentIndex],
			callExpression.Args[argumentIndex+1],
			parameters,
		)...)
	}

	return violations
}

func checkStructuredLogPair(
	fileSet *token.FileSet,
	keyExpression ast.Expr,
	valueExpression ast.Expr,
	parameters policy.GoParameterConfig,
) (violations []Violation) {
	key, found := literalString(keyExpression)
	if !found {
		return []Violation{{
			Position: fileSet.Position(keyExpression.Pos()),
			Rule:     DiagnosticStableLogKeys,
			Message:  "structured log keys must be string literals",
		}}
	}

	if !isStructuredLogKey(key) {
		violations = append(violations, Violation{
			Position: fileSet.Position(keyExpression.Pos()),
			Rule:     DiagnosticStableLogKeys,
			Message: fmt.Sprintf(
				"structured log key %q must be lower-case ASCII with underscores only",
				key,
			),
		})
	}

	if containsSecretLikeName(key, parameters.SecretNames) ||
		expressionContainsSecretLikeIdentifier(valueExpression, parameters.SecretNames) ||
		expressionContainsSensitiveStringLiteral(valueExpression) {
		violations = append(violations, Violation{
			Position: fileSet.Position(valueExpression.Pos()),
			Rule:     DiagnosticNoSecretsInLogs,
			Message:  "structured logs must not include secrets, tokens, or private keys",
		})
	}

	return violations
}

func isStructuredLogCall(selector *ast.SelectorExpr, slogAliases map[string]bool) (found bool) {
	switch selector.Sel.Name {
	case "Debug", "Info", "Warn", "Error":
	default:
		return false
	}

	receiverName, found := rightmostName(selector.X)
	if !found {
		return false
	}

	return slogAliases[receiverName] || receiverName == "logger"
}

func isStructuredLogKey(value string) (valid bool) {
	if value == "" {
		return false
	}

	for _, character := range value {
		if character == '_' ||
			('a' <= character && character <= 'z') ||
			('0' <= character && character <= '9') {
			continue
		}

		return false
	}

	return true
}
