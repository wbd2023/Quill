package syntax

import (
	"go/ast"
	"strings"
	"unicode"
	"unicode/utf8"
)

func isSentinelErrorName(name string) (found bool) {
	if !strings.HasPrefix(name, "Err") || len(name) <= len("Err") {
		return false
	}

	firstSuffixRune, _ := utf8.DecodeRuneInString(name[len("Err"):])
	return unicode.IsUpper(firstSuffixRune)
}

func expressionContainsSecretLikeIdentifier(
	expression ast.Expr,
	secretNames []string,
) (found bool) {
	ast.Inspect(expression, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.Ident:
			if containsSecretLikeName(typed.Name, secretNames) {
				found = true
				return false
			}
		case *ast.SelectorExpr:
			if containsSecretLikeName(typed.Sel.Name, secretNames) {
				found = true
				return false
			}
		}

		return true
	})

	return found
}

func containsSecretLikeName(name string, secretNames []string) (found bool) {
	normalised := strings.ToLower(name)
	for _, secretName := range secretNames {
		fragment := strings.ToLower(secretName)
		if fragment == "" {
			continue
		}

		if strings.Contains(normalised, fragment) {
			return true
		}
	}

	return false
}
