package syntax

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strings"
	"unicode"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
	"github.com/wbd2023/Quill/internal/checks/gopolicy"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// parameter_order constants.
const (
	groupUnknown = -1
)

const minParamFieldSpan = 2

/* --------------------------------------- Ordering Rules --------------------------------------- */

// CheckTypeElision ensures each parameter has its own type declaration.
func CheckTypeElision(fileSet *token.FileSet, file *ast.File) (violations []analysis.Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		funcType, ok := node.(*ast.FuncType)
		if !ok {
			return true
		}

		if funcType.Params == nil {
			return true
		}

		for _, field := range funcType.Params.List {
			if len(field.Names) > 1 {
				names := make([]string, len(field.Names))
				for index, name := range field.Names {
					names[index] = name.Name
				}
				violations = append(violations, analysis.Violation{
					Position: fileSet.Position(field.Pos()),
					Rule:     analysis.DiagnosticNoTypeElision,
					Message: fmt.Sprintf(
						"type elision: parameters %s share a type",
						strings.Join(names, ", "),
					),
				})
			}
		}

		return true
	})

	return violations
}

// CheckParameterOrder ensures ctx is first and secrets are last.
func CheckParameterOrder(
	fileSet *token.FileSet,
	file *ast.File,
	parameters gopolicy.ParameterConfig,
) (violations []analysis.Violation) {
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
					violations = append(violations, analysis.Violation{
						Position: fileSet.Position(name.Pos()),
						Rule:     analysis.DiagnosticContextFirst,
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
				if isSecretParameterName(name.Name, parameters) {
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
			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(funcDecl.Pos()),
				Rule:     analysis.DiagnosticSecretsLast,
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
	constructors gopolicy.ConstructorConfig,
	parameters gopolicy.ParameterConfig,
) (violations []analysis.Violation) {
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

		previousGroup := groupUnknown
		for _, field := range params {
			group := classifyParameter(field, constructors, parameters)
			if group == groupUnknown {
				continue
			}
			if previousGroup != groupUnknown && group < previousGroup {
				violations = append(violations, analysis.Violation{
					Position: fileSet.Position(field.Pos()),
					Rule:     analysis.DiagnosticConstructorCategoryOrder,
					Message: fmt.Sprintf(
						"%s parameter appears after %s parameter in constructor %q",
						parameterGroupName(group, constructors),
						parameterGroupName(previousGroup, constructors),
						funcDecl.Name.Name,
					),
				})
			}
			if group > previousGroup {
				previousGroup = group
			}
		}

		return true
	})
	return violations
}

/* ---------------------------------- Parameter Classification ---------------------------------- */

// isSecretParameterName reports whether name contains a configured secret-name fragment.
func isSecretParameterName(name string, parameters gopolicy.ParameterConfig) (found bool) {
	return containsSecretLikeName(name, parameters.SecretNames)
}

// isConstructor returns true if the function name follows the NewXxx pattern.
func isConstructor(name string) (found bool) {
	if name == "New" {
		return true
	}
	return strings.HasPrefix(name, "New") && len(name) > 3 && unicode.IsUpper(rune(name[3]))
}

// classifyParameter determines the constructor-order group for a parameter.
func classifyParameter(
	field *ast.Field,
	constructors gopolicy.ConstructorConfig,
	parameters gopolicy.ParameterConfig,
) (group int) {
	typeName := typeString(field.Type)

	for index, group := range constructors.ParameterOrder {
		if matchesGroup(field, typeName, group, parameters) {
			return index
		}
	}

	return groupUnknown
}

func matchesGroup(
	field *ast.Field,
	typeName string,
	group gopolicy.ParameterGroup,
	parameters gopolicy.ParameterConfig,
) (matches bool) {
	if hasAnySuffix(typeName, group.TypeNameSuffixes) {
		return true
	}

	for _, name := range field.Names {
		if slices.Contains(group.ParameterNames, name.Name) {
			return true
		}

		if group.MatchesSecretNames && isSecretParameterName(name.Name, parameters) {
			return true
		}
	}

	return false
}

func parameterGroupName(
	group int,
	constructors gopolicy.ConstructorConfig,
) (name string) {
	if group < 0 || group >= len(constructors.ParameterOrder) {
		return "unknown"
	}

	return constructors.ParameterOrder[group].Name
}

func hasAnySuffix(target string, suffixes []string) (found bool) {
	return slices.ContainsFunc(suffixes, func(suffix string) bool {
		return strings.HasSuffix(target, suffix)
	})
}
