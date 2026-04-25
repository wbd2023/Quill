package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"ciphera/tools/internal/profile"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	categoryUnknown = -1
)

const minParamFieldSpan = 2

/* --------------------------------------- Ordering Rules --------------------------------------- */

// CheckParameterOrder ensures ctx is first and secrets are last.
func CheckParameterOrder(
	fileSet *token.FileSet,
	file *ast.File,
	parameters profile.GoParameterConfig,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok || funcDecl.Type.Params == nil {
			return true
		}

		params := funcDecl.Type.Params.List
		if len(params) < minParamFieldSpan {
			return true
		}

		funcName := funcDecl.Name.Name

		for fieldIndex, field := range params {
			for _, name := range field.Names {
				if name.Name == "ctx" && fieldIndex > 0 {
					violations = append(violations, Violation{
						Position: fileSet.Position(name.Pos()),
						Rule:     DiagnosticContextFirst,
						Message:  fmt.Sprintf("ctx must be the first parameter in %q", funcName),
					})
				}
			}
		}

		lastNonSecretIndex := -1
		firstSecretIndex := len(params)
		for fieldIndex, field := range params {
			isSecret := false
			for _, name := range field.Names {
				if isSecretName(name.Name, parameters) {
					isSecret = true
				}
			}
			if isSecret && fieldIndex < firstSecretIndex {
				firstSecretIndex = fieldIndex
			}
			if !isSecret {
				lastNonSecretIndex = fieldIndex
			}
		}
		if firstSecretIndex < lastNonSecretIndex {
			violations = append(violations, Violation{
				Position: fileSet.Position(funcDecl.Pos()),
				Rule:     DiagnosticSecretsLast,
				Message:  fmt.Sprintf("secret parameters must be last in %q", funcName),
			})
		}

		return true
	})
	return violations
}

// CheckConstructorOrder ensures constructor parameters follow the profile's canonical order.
func CheckConstructorOrder(
	fileSet *token.FileSet,
	file *ast.File,
	parameters profile.GoParameterConfig,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok || funcDecl.Type.Params == nil {
			return true
		}

		if !isConstructor(funcDecl.Name.Name) {
			return true
		}

		params := funcDecl.Type.Params.List
		if len(params) < minParamFieldSpan {
			return true
		}

		prevCategory := categoryUnknown
		for _, field := range params {
			category := classifyParam(field, parameters)
			if category == categoryUnknown {
				continue
			}
			if prevCategory != categoryUnknown && category < prevCategory {
				violations = append(violations, Violation{
					Position: fileSet.Position(field.Pos()),
					Rule:     DiagnosticConstructorCategoryOrder,
					Message: fmt.Sprintf(
						"%s parameter appears after %s parameter in constructor %q",
						parameterCategoryLabel(category, parameters),
						parameterCategoryLabel(prevCategory, parameters),
						funcDecl.Name.Name,
					),
				})
			}
			if category > prevCategory {
				prevCategory = category
			}
		}

		return true
	})
	return violations
}

/* ---------------------------------- Parameter Classification ---------------------------------- */

// isSecretName returns true if the parameter name represents a secret.
func isSecretName(name string, parameters profile.GoParameterConfig) (found bool) {
	return contains(parameters.SecretNames, name)
}

// isConstructor returns true if the function name follows the NewXxx pattern.
func isConstructor(name string) (found bool) {
	if name == "New" {
		return true
	}
	return strings.HasPrefix(name, "New") && len(name) > 3 && unicode.IsUpper(rune(name[3]))
}

// classifyParam determines the category of a constructor parameter.
func classifyParam(field *ast.Field, parameters profile.GoParameterConfig) (category int) {
	typeName := typeString(field.Type)

	for index, current := range parameters.ConstructorCategories {
		if matchesCategory(field, typeName, current, parameters) {
			return index
		}
	}

	return categoryUnknown
}

func matchesCategory(
	field *ast.Field,
	typeName string,
	category profile.GoConstructorCategory,
	parameters profile.GoParameterConfig,
) (matches bool) {
	if containsAny(typeName, category.ExcludedTypeMarkers) {
		return false
	}

	if containsAny(typeName, category.TypeMarkers) {
		return true
	}

	for _, name := range field.Names {
		if contains(category.ParameterNames, name.Name) {
			return true
		}

		if category.UsesSecretNames && isSecretName(name.Name, parameters) {
			return true
		}
	}

	return false
}

func parameterCategoryLabel(
	category int,
	parameters profile.GoParameterConfig,
) (label string) {
	if category < 0 || category >= len(parameters.ConstructorCategories) {
		return "unknown"
	}

	return parameters.ConstructorCategories[category].Name
}

func contains(values []string, target string) (found bool) {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func containsAny(target string, fragments []string) (found bool) {
	for _, fragment := range fragments {
		if strings.Contains(target, fragment) {
			return true
		}
	}

	return false
}
