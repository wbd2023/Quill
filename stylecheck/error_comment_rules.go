package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const adaptersPathSegment = "/internal/adapters/"
const cmdPathSegment = "/cmd/"
const internalPathSegment = "/internal/"
const testsPathSegment = "/tests/"
const domainErrorsFilePathSuffix = "/internal/core/domain/errors.go"
const inlineCommentDirectiveCodeGenerated = "code generated"
const inlineCommentDirectiveFixme = "fixme:"
const inlineCommentDirectiveGo = "go:"
const inlineCommentDirectiveNolint = "nolint"
const inlineCommentDirectiveTodo = "todo:"
const inlineCommentPunctuation = ".!?"
const secretLikeNameFragmentPassphrase = "passphrase"
const secretLikeNameFragmentPassword = "password"
const secretLikeNameFragmentPrivateKey = "privatekey"
const secretLikeNameFragmentSecretKey = "secretkey"
const secretLikeNameFragmentSecret = "secret"
const secretLikeNameFragmentToken = "token"
const secretLikeNameFragmentSeed = "seed"

/* ----------------------------------------- Error Rules ---------------------------------------- */

// checkGoErrorHandlingStyle enforces Go error-message and sentinel-error style (2.1).
func checkGoErrorHandlingStyle(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
) (violations []violation) {
	if !isAppScopePath(path) {
		return nil
	}

	fmtImportAliases := importAliasesForPath(file, "fmt")
	errorsImportAliases := importAliasesForPath(file, "errors")

	ast.Inspect(file, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok || len(callExpression.Args) == 0 {
			return true
		}

		selectorExpression, ok := callExpression.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		packageIdentifier, ok := selectorExpression.X.(*ast.Ident)
		if !ok {
			return true
		}

		switch {
		case selectorExpression.Sel.Name == "Errorf" && fmtImportAliases[packageIdentifier.Name]:
			message, found := extractStringLiteral(callExpression.Args[0])
			if found {
				violations = append(
					violations,
					checkErrorMessageLiteralStyle(
						fileSet,
						callExpression.Args[0],
						message,
						"fmt.Errorf",
					)...,
				)
			}

			if isTestFile {
				return true
			}

			for _, arg := range callExpression.Args[1:] {
				if !expressionContainsSecretLikeIdentifier(arg) {
					continue
				}

				violations = append(violations, violation{
					position: fileSet.Position(arg.Pos()),
					rule:     "2.1",
					message:  "error context must not include secrets in fmt.Errorf arguments",
				})
			}

		case selectorExpression.Sel.Name == "New" && errorsImportAliases[packageIdentifier.Name]:
			message, found := extractStringLiteral(callExpression.Args[0])
			if !found {
				return true
			}

			violations = append(
				violations,
				checkErrorMessageLiteralStyle(
					fileSet,
					callExpression.Args[0],
					message,
					"errors.New",
				)...,
			)
		}

		return true
	})

	if isTestFile || strings.HasSuffix(path, domainErrorsFilePathSuffix) {
		return violations
	}

	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, specification := range genDecl.Specs {
			valueSpec, ok := specification.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, name := range valueSpec.Names {
				if !isSentinelErrorName(name.Name) {
					continue
				}

				violations = append(violations, violation{
					position: fileSet.Position(name.Pos()),
					rule:     "2.1",
					message:  "sentinel errors must be declared in internal/core/domain/errors.go",
				})
			}
		}
	}

	return violations
}

// checkAdapterErrorWrapping rejects bare error propagation in adapters (2.1).
func checkAdapterErrorWrapping(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
) (violations []violation) {
	if isTestFile || !strings.Contains(path, adaptersPathSegment) {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		functionDecl, ok := node.(*ast.FuncDecl)
		if !ok || functionDecl.Body == nil {
			return true
		}

		ast.Inspect(functionDecl.Body, func(bodyNode ast.Node) bool {
			switch typed := bodyNode.(type) {
			case *ast.FuncLit:
				return false
			case *ast.ReturnStmt:
				if !isBareErrReturn(typed) {
					return true
				}

				violations = append(violations, violation{
					position: fileSet.Position(typed.Return),
					rule:     "2.1",
					message:  "adapter error returns must wrap low-level errors with context (%w)",
				})
			}

			return true
		})

		return false
	})

	return violations
}

/* ---------------------------------------- Comment Rules --------------------------------------- */

// checkInlineCommentStyle validates trailing inline comment case and punctuation (2.3).
func checkInlineCommentStyle(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if !isAppScopePath(path) {
		return nil
	}

	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)
	seen := make(map[token.Pos]bool)

	for node, commentGroups := range commentMap {
		nodeEndLine := fileSet.Position(node.End()).Line

		for _, commentGroup := range commentGroups {
			for _, comment := range commentGroup.List {
				if !strings.HasPrefix(comment.Text, "//") {
					continue
				}

				if seen[comment.Pos()] {
					continue
				}

				commentPosition := fileSet.Position(comment.Pos())
				if commentPosition.Line != nodeEndLine {
					continue
				}

				payload := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
				if payload == "" || isInlineCommentDirective(payload) {
					continue
				}

				seen[comment.Pos()] = true

				if startsWithUppercaseLetter(payload) {
					violations = append(violations, violation{
						position: fileSet.Position(comment.Pos()),
						rule:     "2.3",
						message:  "inline trailing comment should start lower-case",
					})
				}

				if endsWithSentencePunctuation(payload) {
					violations = append(violations, violation{
						position: fileSet.Position(comment.Pos()),
						rule:     "2.3",
						message:  "inline trailing comment should not end with punctuation",
					})
				}
			}
		}
	}

	return violations
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func isAppScopePath(path string) (found bool) {
	return strings.Contains(path, internalPathSegment) ||
		strings.Contains(path, cmdPathSegment) ||
		strings.Contains(path, testsPathSegment)
}

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

func extractStringLiteral(expression ast.Expr) (value string, found bool) {
	literal, ok := expression.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return "", false
	}

	unquotedValue, err := strconv.Unquote(literal.Value)
	if err != nil {
		return "", false
	}

	return unquotedValue, true
}

func checkErrorMessageLiteralStyle(
	fileSet *token.FileSet,
	expression ast.Expr,
	message string,
	callName string,
) (violations []violation) {
	trimmedMessage := strings.TrimSpace(message)
	if trimmedMessage == "" {
		return nil
	}

	if startsWithUppercaseLetter(trimmedMessage) {
		violations = append(violations, violation{
			position: fileSet.Position(expression.Pos()),
			rule:     "2.1",
			message:  fmt.Sprintf("error context must be lowercase (%s)", callName),
		})
	}

	if endsWithSentencePunctuation(trimmedMessage) {
		violations = append(violations, violation{
			position: fileSet.Position(expression.Pos()),
			rule:     "2.1",
			message:  fmt.Sprintf("error context must not end with punctuation (%s)", callName),
		})
	}

	return violations
}

func isSentinelErrorName(name string) (found bool) {
	if !strings.HasPrefix(name, "Err") || len(name) <= len("Err") {
		return false
	}

	firstSuffixRune, _ := utf8.DecodeRuneInString(name[len("Err"):])
	return unicode.IsUpper(firstSuffixRune)
}

func expressionContainsSecretLikeIdentifier(expression ast.Expr) (found bool) {
	ast.Inspect(expression, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.Ident:
			if containsSecretLikeName(typed.Name) {
				found = true
				return false
			}
		case *ast.SelectorExpr:
			if containsSecretLikeName(typed.Sel.Name) {
				found = true
				return false
			}
		}

		return true
	})

	return found
}

func containsSecretLikeName(name string) (found bool) {
	normalised := strings.ToLower(name)

	return strings.Contains(normalised, secretLikeNameFragmentPassphrase) ||
		strings.Contains(normalised, secretLikeNameFragmentPassword) ||
		strings.Contains(normalised, secretLikeNameFragmentPrivateKey) ||
		strings.Contains(normalised, secretLikeNameFragmentSecretKey) ||
		strings.Contains(normalised, secretLikeNameFragmentSecret) ||
		strings.Contains(normalised, secretLikeNameFragmentToken) ||
		strings.Contains(normalised, secretLikeNameFragmentSeed)
}

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

func isInlineCommentDirective(comment string) (found bool) {
	normalisedComment := strings.ToLower(strings.TrimSpace(comment))

	return strings.HasPrefix(normalisedComment, inlineCommentDirectiveNolint) ||
		strings.HasPrefix(normalisedComment, inlineCommentDirectiveTodo) ||
		strings.HasPrefix(normalisedComment, inlineCommentDirectiveFixme) ||
		strings.HasPrefix(normalisedComment, inlineCommentDirectiveGo) ||
		strings.HasPrefix(normalisedComment, inlineCommentDirectiveCodeGenerated)
}

func startsWithUppercaseLetter(value string) (found bool) {
	firstRune, _ := utf8.DecodeRuneInString(value)
	return unicode.IsUpper(firstRune)
}

func endsWithSentencePunctuation(value string) (found bool) {
	lastRune, _ := utf8.DecodeLastRuneInString(value)
	return strings.ContainsRune(inlineCommentPunctuation, lastRune)
}
