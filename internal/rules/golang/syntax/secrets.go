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
			violations = append(violations, secretValueSpecViolations(
				fileSet,
				typedNode,
				parameters.SecretNames,
			)...)

		case *ast.AssignStmt:
			violations = append(violations, secretAssignmentViolations(
				fileSet,
				typedNode,
				parameters.SecretNames,
			)...)

		case *ast.KeyValueExpr:
			violation, found := secretFieldViolation(
				fileSet,
				typedNode,
				parameters.SecretNames,
			)
			if found {
				violations = append(violations, violation)
			}
		}

		return true
	})

	return violations
}

/* --------------------------------------- AST Node Scans --------------------------------------- */

func secretValueSpecViolations(
	fileSet *token.FileSet,
	valueSpec *ast.ValueSpec,
	secretNames []string,
) (violations []analysis.Violation) {
	for index, name := range valueSpec.Names {
		if index >= len(valueSpec.Values) || !containsSecretLikeName(name.Name, secretNames) {
			continue
		}

		literal, found := secretStringLiteral(valueSpec.Values[index])
		if !found {
			continue
		}

		violations = append(violations, secretLiteralViolation(fileSet, literal))
	}

	return violations
}

func secretAssignmentViolations(
	fileSet *token.FileSet,
	assignment *ast.AssignStmt,
	secretNames []string,
) (violations []analysis.Violation) {
	for index, lhsExpression := range assignment.Lhs {
		if index >= len(assignment.Rhs) {
			continue
		}

		name, found := rightmostName(lhsExpression)
		if !found || !containsSecretLikeName(name, secretNames) {
			continue
		}

		literal, found := secretStringLiteral(assignment.Rhs[index])
		if !found {
			continue
		}

		violations = append(violations, secretLiteralViolation(fileSet, literal))
	}

	return violations
}

func secretFieldViolation(
	fileSet *token.FileSet,
	field *ast.KeyValueExpr,
	secretNames []string,
) (violation analysis.Violation, found bool) {
	keyName, found := rightmostName(field.Key)
	if !found || !containsSecretLikeName(keyName, secretNames) {
		return analysis.Violation{}, false
	}

	literal, found := secretStringLiteral(field.Value)
	if !found {
		return analysis.Violation{}, false
	}

	return secretLiteralViolation(fileSet, literal), true
}

func secretLiteralViolation(
	fileSet *token.FileSet,
	expression ast.Expr,
) (violation analysis.Violation) {
	return analysis.Violation{
		Position: fileSet.Position(expression.Pos()),
		Rule:     analysis.DiagnosticNoSecretsInRepository,
		Message:  "source code must not hard-code secret-like string literals",
	}
}

/* -------------------------------------- Expression Scans -------------------------------------- */

func expressionContainsSensitiveStringLiteral(expression ast.Expr) (found bool) {
	_, found = secretStringLiteral(expression)
	return found
}

func secretStringLiteral(expression ast.Expr) (literal ast.Expr, found bool) {
	ast.Inspect(expression, func(node ast.Node) bool {
		basicLiteral, ok := node.(*ast.BasicLit)
		if !ok || basicLiteral.Kind != token.STRING {
			return true
		}

		value, ok := literalString(basicLiteral)
		if !ok || !looksSecretLiteral(value) {
			return true
		}

		literal = basicLiteral
		found = true
		return false
	})

	return literal, found
}

/* ----------------------------------- Literal Classification ----------------------------------- */

func looksSecretLiteral(value string) (found bool) {
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
