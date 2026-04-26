package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"ciphera/tools/internal/policy"
)

/* ----------------------------------------- Error Rules ---------------------------------------- */

// CheckErrorHandlingStyle enforces Go error-message and sentinel-error style.
func CheckErrorHandlingStyle(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
	parameters policy.GoParameterConfig,
) (violations []Violation) {
	if !classifier.HasClass(path, PathClassGoSource) {
		return nil
	}

	fmtImportAliases := importAliasesForPath(file, "fmt")
	errorsImportAliases := importAliasesForPath(file, "errors")

	ast.Inspect(file, func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok || len(callExpr.Args) == 0 {
			return true
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkgIdent, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			return true
		}

		switch {
		case selectorExpr.Sel.Name == "Errorf" && fmtImportAliases[pkgIdent.Name]:
			message, found := literalString(callExpr.Args[0])
			if found {
				violations = append(
					violations,
					checkErrorMessageLiteralStyle(
						fileSet,
						callExpr.Args[0],
						message,
						"fmt.Errorf",
					)...,
				)
			}

			if isTestFile {
				return true
			}

			for _, arg := range callExpr.Args[1:] {
				if !expressionContainsSecretLikeIdentifier(arg, parameters.SecretNames) {
					continue
				}

				violations = append(violations, Violation{
					Position: fileSet.Position(arg.Pos()),
					Rule:     DiagnosticErrorContextNoSecrets,
					Message:  "error context must not include secrets in fmt.Errorf arguments",
				})
			}

		case selectorExpr.Sel.Name == "New" && errorsImportAliases[pkgIdent.Name]:
			message, found := literalString(callExpr.Args[0])
			if !found {
				return true
			}

			violations = append(
				violations,
				checkErrorMessageLiteralStyle(
					fileSet,
					callExpr.Args[0],
					message,
					"errors.New",
				)...,
			)
		}

		return true
	})

	if isTestFile || classifier.HasClass(path, PathClassDomainErrors) {
		return violations
	}

	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, name := range valueSpec.Names {
				if !isSentinelErrorName(name.Name) {
					continue
				}

				violations = append(violations, Violation{
					Position: fileSet.Position(name.Pos()),
					Rule:     DiagnosticDomainErrorsLocation,
					Message: fmt.Sprintf(
						"sentinel errors must be declared in %s",
						classifier.FirstPattern(PathClassDomainErrors),
					),
				})
			}
		}
	}

	return violations
}

// CheckAdapterErrorWrapping rejects bare error propagation in adapters.
func CheckAdapterErrorWrapping(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
) (violations []Violation) {
	if isTestFile || !classifier.HasClass(path, PathClassConcreteInfra) {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		ast.Inspect(funcDecl.Body, func(bodyNode ast.Node) bool {
			switch typed := bodyNode.(type) {
			case *ast.FuncLit:
				return false

			case *ast.ReturnStmt:
				if !isBareErrReturn(typed) {
					return true
				}

				violations = append(violations, Violation{
					Position: fileSet.Position(typed.Return),
					Rule:     DiagnosticAdapterWrapsCause,
					Message:  "adapter error returns must wrap low-level errors with context (%w)",
				})
			}

			return true
		})

		return false
	})

	return violations
}

/* --------------------------------------- Import Aliases --------------------------------------- */

func importAliasesForPath(file *ast.File, importPath string) (aliases map[string]bool) {
	aliases = make(map[string]bool)

	for _, importSpec := range file.Imports {
		if importSpec.Path == nil || importSpec.Path.Kind != token.STRING {
			continue
		}

		importedPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil || importedPath != importPath {
			continue
		}

		if importSpec.Name == nil {
			aliases[pathBase(importPath)] = true
			continue
		}

		if importSpec.Name.Name == "." || importSpec.Name.Name == "_" {
			continue
		}

		aliases[importSpec.Name.Name] = true
	}

	return aliases
}

func pathBase(value string) (base string) {
	if value == "" {
		return ""
	}

	parts := strings.Split(value, "/")
	return parts[len(parts)-1]
}

/* ---------------------------------------- Message Rules --------------------------------------- */

func checkErrorMessageLiteralStyle(
	fileSet *token.FileSet,
	expression ast.Expr,
	message string,
	callName string,
) (violations []Violation) {
	trimmedMessage := strings.TrimSpace(message)
	if trimmedMessage == "" {
		return nil
	}

	if startsWithUppercaseLetter(trimmedMessage) {
		violations = append(violations, Violation{
			Position: fileSet.Position(expression.Pos()),
			Rule:     DiagnosticErrorContextLowercase,
			Message:  fmt.Sprintf("error context must be lowercase (%s)", callName),
		})
	}

	if endsWithSentencePunctuation(trimmedMessage) {
		violations = append(violations, Violation{
			Position: fileSet.Position(expression.Pos()),
			Rule:     DiagnosticErrorContextNoPunctuation,
			Message:  fmt.Sprintf("error context must not end with punctuation (%s)", callName),
		})
	}

	return violations
}

/* -------------------------------------- Secret Detection -------------------------------------- */

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

/* ---------------------------------------- Return Rules ---------------------------------------- */

func isBareErrReturn(returnStatement *ast.ReturnStmt) (found bool) {
	if len(returnStatement.Results) == 0 {
		return false
	}

	lastReturnExpression := returnStatement.Results[len(returnStatement.Results)-1]
	identifier, ok := lastReturnExpression.(*ast.Ident)
	if !ok {
		return false
	}

	return identifier.Name == "err"
}
