package checker

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"stylecheck/internal/checker/support"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// Constructor parameter category ordering.
const (
	categoryUnknown    = 0
	categoryRepository = 1
	categoryService    = 2
	categoryAdapter    = 3
	categoryConfig     = 4
	categorySecret     = 5
)

const minParamFieldSpan = 2

var constructorParameterCategoryLabels = map[int]string{
	categoryRepository: "repository",
	categoryService:    "service",
	categoryAdapter:    "adapter",
	categoryConfig:     "config",
	categorySecret:     "secret",
}

var secretParameterNames = map[string]bool{
	"passphrase": true,
	"privateKey": true,
	"token":      true,
	"seed":       true,
	"secret":     true,
	"password":   true,
	"secretKey":  true,
}

var configParameterNames = map[string]bool{
	"serverURL":  true,
	"relayURL":   true,
	"identityID": true,
	"timeout":    true,
}

/* --------------------------------------- Ordering Rules --------------------------------------- */

// checkParamOrder ensures ctx is first and secrets are last (2.7).
func checkParamOrder(fileSet *token.FileSet, file *ast.File) (violations []violation) {
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
					violations = append(violations, violation{
						position: fileSet.Position(name.Pos()),
						rule:     "2.7",
						message:  fmt.Sprintf("ctx must be the first parameter in %q", funcName),
					})
				}
			}
		}

		lastNonSecretIndex := -1
		firstSecretIndex := len(params)
		for fieldIndex, field := range params {
			isSecret := false
			for _, name := range field.Names {
				if isSecretName(name.Name) {
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
			violations = append(violations, violation{
				position: fileSet.Position(funcDecl.Pos()),
				rule:     "2.7",
				message:  fmt.Sprintf("secret parameters must be last in %q", funcName),
			})
		}

		return true
	})
	return violations
}

// checkConstructorOrder ensures constructor parameters follow the canonical ordering:
// repositories -> services -> adapters -> config -> secrets (2.8).
func checkConstructorOrder(fileSet *token.FileSet, file *ast.File) (violations []violation) {
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
			category := classifyParam(field)
			if category == categoryUnknown {
				continue
			}
			if prevCategory != categoryUnknown && category < prevCategory {
				violations = append(violations, violation{
					position: fileSet.Position(field.Pos()),
					rule:     "2.8",
					message: fmt.Sprintf(
						"%s parameter appears after %s parameter in constructor %q",
						constructorParameterCategoryLabels[category],
						constructorParameterCategoryLabels[prevCategory],
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

/* ----------------------------------- Classification Helpers ----------------------------------- */

// isSecretName returns true if the parameter name represents a secret.
func isSecretName(name string) (found bool) {
	return secretParameterNames[name]
}

// isConfigName returns true if the parameter name represents configuration.
func isConfigName(name string) (found bool) {
	return configParameterNames[name]
}

// isConstructor returns true if the function name follows the NewXxx pattern.
func isConstructor(name string) (found bool) {
	if name == "New" {
		return true
	}
	return strings.HasPrefix(name, "New") && len(name) > 3 && unicode.IsUpper(rune(name[3]))
}

// classifyParam determines the category of a constructor parameter.
func classifyParam(field *ast.Field) (category int) {
	typeName := support.TypeString(field.Type)

	if strings.Contains(typeName, "Repository") {
		return categoryRepository
	}

	if strings.Contains(typeName, "Service") && !strings.Contains(typeName, "Config") {
		return categoryService
	}

	if strings.Contains(typeName, "Client") || strings.Contains(typeName, "Factory") {
		return categoryAdapter
	}

	for _, name := range field.Names {
		if isSecretName(name.Name) {
			return categorySecret
		}

		if isConfigName(name.Name) {
			return categoryConfig
		}
	}

	return categoryUnknown
}
