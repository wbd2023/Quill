package checks

import (
	"go/ast"
	"go/token"
	"strings"

	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const minSensitiveLiteralLength = 16

/* ------------------------------------ Sensitive Data Rules ------------------------------------ */

func CheckSensitiveDataLiterals(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
	parameters profile.GoParameterConfig,
) (violations []Violation) {
	if !classifier.HasClass(path, rulepack.PathClassApp) ||
		isTestFile ||
		classifier.HasClass(path, rulepack.PathClassTestMocks) {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		switch typedNode := node.(type) {
		case *ast.ValueSpec:
			for index, name := range typedNode.Names {
				if !containsSecretLikeName(name.Name, parameters.SecretNames) ||
					index >= len(typedNode.Values) {
					continue
				}

				if !expressionContainsSensitiveStringLiteral(typedNode.Values[index]) {
					continue
				}

				violations = append(violations, Violation{
					Position: fileSet.Position(typedNode.Values[index].Pos()),
					Rule:     DiagnosticNoSecretsInRepository,
					Message:  "source code must not hard-code secret-like string literals",
				})
			}

		case *ast.AssignStmt:
			for index, lhsExpression := range typedNode.Lhs {
				if index >= len(typedNode.Rhs) {
					continue
				}

				name, found := assignedName(lhsExpression)
				if !found || !containsSecretLikeName(name, parameters.SecretNames) {
					continue
				}
				if !expressionContainsSensitiveStringLiteral(typedNode.Rhs[index]) {
					continue
				}

				violations = append(violations, Violation{
					Position: fileSet.Position(typedNode.Rhs[index].Pos()),
					Rule:     DiagnosticNoSecretsInRepository,
					Message:  "source code must not hard-code secret-like string literals",
				})
			}

		case *ast.KeyValueExpr:
			keyName, found := keyedFieldName(typedNode.Key)
			if !found || !containsSecretLikeName(keyName, parameters.SecretNames) {
				return true
			}
			if !expressionContainsSensitiveStringLiteral(typedNode.Value) {
				return true
			}

			violations = append(violations, Violation{
				Position: fileSet.Position(typedNode.Value.Pos()),
				Rule:     DiagnosticNoSecretsInRepository,
				Message:  "source code must not hard-code secret-like string literals",
			})
		}

		return true
	})

	return violations
}

/* -------------------------------------- Expression Scans -------------------------------------- */

func expressionContainsSensitiveStringLiteral(expression ast.Expr) (found bool) {
	ast.Inspect(expression, func(node ast.Node) bool {
		literal, ok := node.(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return true
		}

		value, ok := stringLiteralValue(literal)
		if !ok || !looksSensitiveLiteral(value) {
			return true
		}

		found = true
		return false
	})

	return found
}

/* ----------------------------------- Literal Classification ----------------------------------- */

func looksSensitiveLiteral(value string) (found bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}

	if strings.HasPrefix(trimmed, "-----BEGIN ") {
		return true
	}

	if isLongHexLikeString(trimmed) || isLongBase64LikeString(trimmed) {
		return true
	}

	return len(trimmed) >= minSensitiveLiteralLength &&
		containsLetter(trimmed) &&
		containsDigit(trimmed)
}

/* --------------------------------------- Name Extraction -------------------------------------- */

func assignedName(expression ast.Expr) (name string, found bool) {
	switch typedExpression := expression.(type) {
	case *ast.Ident:
		return typedExpression.Name, true
	case *ast.SelectorExpr:
		return typedExpression.Sel.Name, true
	default:
		return "", false
	}
}

func keyedFieldName(expression ast.Expr) (name string, found bool) {
	switch typedExpression := expression.(type) {
	case *ast.Ident:
		return typedExpression.Name, true
	case *ast.SelectorExpr:
		return typedExpression.Sel.Name, true
	default:
		return "", false
	}
}

/* ------------------------------------ String Classification ----------------------------------- */

func isLongHexLikeString(value string) (found bool) {
	if len(value) < 32 || len(value)%2 != 0 {
		return false
	}

	for _, runeValue := range value {
		switch {
		case runeValue >= '0' && runeValue <= '9':

		case runeValue >= 'a' && runeValue <= 'f':

		case runeValue >= 'A' && runeValue <= 'F':

		default:
			return false
		}
	}

	return true
}

func isLongBase64LikeString(value string) (found bool) {
	if len(value) < 24 || len(value)%4 != 0 {
		return false
	}

	for _, runeValue := range value {
		switch {
		case runeValue >= 'a' && runeValue <= 'z':

		case runeValue >= 'A' && runeValue <= 'Z':

		case runeValue >= '0' && runeValue <= '9':

		case runeValue == '+' || runeValue == '/' || runeValue == '=':

		default:
			return false
		}
	}

	return true
}

func containsLetter(value string) (found bool) {
	for _, runeValue := range value {
		if (runeValue >= 'a' && runeValue <= 'z') || (runeValue >= 'A' && runeValue <= 'Z') {
			return true
		}
	}

	return false
}

func containsDigit(value string) (found bool) {
	for _, runeValue := range value {
		if runeValue >= '0' && runeValue <= '9' {
			return true
		}
	}

	return false
}
