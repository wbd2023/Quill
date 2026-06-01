package syntax

import (
	"go/ast"
	"go/token"
	"strings"

	"ciphera/tools/internal/rules/golang/analysis"
	gopolicy "ciphera/tools/internal/rules/golang/policy"
)

const minSensitiveLiteralLength = 16

/* ------------------------------------ Sensitive Data Rules ------------------------------------ */

func CheckSensitiveDataLiterals(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier analysis.PathClassifier,
	parameters gopolicy.ParameterConfig,
) (violations []analysis.Violation) {
	if !classifier.HasRole(path, analysis.PathRoleGoSource) ||
		isTestFile ||
		classifier.HasRole(path, analysis.PathRoleTestMocks) {
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

				violations = append(violations, analysis.Violation{
					Position: fileSet.Position(typedNode.Values[index].Pos()),
					Rule:     analysis.DiagnosticNoSecretsInRepository,
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

				violations = append(violations, analysis.Violation{
					Position: fileSet.Position(typedNode.Rhs[index].Pos()),
					Rule:     analysis.DiagnosticNoSecretsInRepository,
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

			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(typedNode.Value.Pos()),
				Rule:     analysis.DiagnosticNoSecretsInRepository,
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

		value, ok := literalString(literal)
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

	for _, character := range value {
		switch {
		case character >= '0' && character <= '9':

		case character >= 'a' && character <= 'f':

		case character >= 'A' && character <= 'F':

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

	for _, character := range value {
		switch {
		case character >= 'a' && character <= 'z':

		case character >= 'A' && character <= 'Z':

		case character >= '0' && character <= '9':

		case character == '+' || character == '/' || character == '=':

		default:
			return false
		}
	}

	return true
}

func containsLetter(value string) (found bool) {
	return strings.ContainsFunc(value, func(character rune) bool {
		return (character >= 'a' && character <= 'z') || (character >= 'A' && character <= 'Z')
	})
}

func containsDigit(value string) (found bool) {
	return strings.ContainsFunc(value, func(character rune) bool {
		return character >= '0' && character <= '9'
	})
}
