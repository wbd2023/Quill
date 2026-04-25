package checks

import (
	"fmt"
	"go/ast"
	"go/token"

	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
)

/* ------------------------------------- Structured Logging ------------------------------------- */

func CheckStructuredLogging(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	classifier PathClassifier,
	parameters profile.GoParameterConfig,
) (violations []Violation) {
	if !classifier.HasClass(path, rulepack.PathClassApp) {
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

		selector, ok := callExpression.Fun.(*ast.SelectorExpr)
		if !ok || !isStructuredLogCall(selector, slogAliases) {
			return true
		}

		fieldCount := len(callExpression.Args) - 1
		if fieldCount <= 0 {
			return true
		}

		if fieldCount%2 != 0 {
			violations = append(violations, Violation{
				Position: fileSet.Position(callExpression.Pos()),
				Rule:     DiagnosticStructuredLogs,
				Message:  "structured log calls must use key/value pairs",
			})
			return true
		}

		for argumentIndex := 1; argumentIndex < len(callExpression.Args); argumentIndex += 2 {
			keyExpression := callExpression.Args[argumentIndex]
			key, found := stringLiteralValue(keyExpression)
			if !found {
				violations = append(violations, Violation{
					Position: fileSet.Position(keyExpression.Pos()),
					Rule:     DiagnosticStableLogKeys,
					Message:  "structured log keys must be string literals",
				})
				continue
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

			valueExpression := callExpression.Args[argumentIndex+1]
			if containsSecretLikeName(key, parameters.SecretNames) ||
				expressionContainsSecretLikeIdentifier(valueExpression, parameters.SecretNames) ||
				expressionContainsSensitiveStringLiteral(valueExpression) {
				violations = append(violations, Violation{
					Position: fileSet.Position(valueExpression.Pos()),
					Rule:     DiagnosticNoSecretsInLogs,
					Message:  "structured logs must not include secrets, tokens, or private keys",
				})
			}
		}

		return true
	})

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

	for _, runeValue := range value {
		if runeValue == '_' ||
			('a' <= runeValue && runeValue <= 'z') ||
			('0' <= runeValue && runeValue <= '9') {
			continue
		}

		return false
	}

	return true
}
