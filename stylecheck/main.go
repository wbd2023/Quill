// Command stylecheck performs AST-based style checks on Go source files.
//
// It enforces the following rules from STYLE.md:
//   - 2.2 Named returns: all functions must use named, descriptive return values.
//   - 2.2 Naked returns: explicit return values are required.
//   - 2.2 Type elision: each parameter must have its own type.
//   - 2.2 Domain ID constructors: avoid direct casts for key domain identifier types.
//     This uses a type-aware pass with syntax fallback for non-buildable snippets.
//   - 2.1 Error handling: lowercase/no-punctuation error context, no secrets in fmt.Errorf args,
//     and sentinel errors scoped to domain/errors.go.
//   - 2.1 Adapter error wrapping: reject bare `return err` propagation in adapters.
//   - 2.3 Inline comment style: trailing comments must start lower-case and avoid punctuation.
//   - 2.2 Single-letter variables: only i, j, k (loops) and receivers.
//   - 2.2 Service package type naming: exported types end with Service/UseCase/Config.
//   - 2.5 CRUD-L ordering inside interfaces.
//   - 2.5 Mock method order matches interface method order exactly.
//   - 2.5 Implementation method order matches interface method order exactly.
//   - 2.7 Parameter ordering: ctx first, secrets last.
//   - 2.8 Constructor ordering: repos -> services -> adapters -> config -> secrets.
//   - 2.9 File structure ordering for top-level declarations.
//
// Usage:
//
//	stylecheck <dir>...
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/packages"
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

const (
	crudUnknown = 0
	crudCreate  = 1
	crudRead    = 2
	crudUpdate  = 3
	crudDelete  = 4
	crudList    = 5
)

const (
	declUnknown    = 0
	declConstants  = 1
	declErrors     = 2
	declTypes      = 3
	declAssertions = 4
)

const (
	minRequiredArgs   = 2
	usageExitCode     = 2
	minParamFieldSpan = 2
	minCategoryKinds  = 2
)

const domainPackagePathSegment = "/internal/core/domain/"
const domainPackagePathSuffix = "/internal/core/domain"
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

/* -------------------------------------------- Types ------------------------------------------- */

// violation represents a single style rule violation.
type violation struct {
	position token.Position
	rule     string
	message  string
}

type methodDecl struct {
	name     string
	position token.Position
}

type interfaceDecl struct {
	name     string
	methods  []methodDecl
	position token.Position
}

type implementationBinding struct {
	interfaceName      string
	implementationName string
	implementationKey  string
	position           token.Position
}

type analysisState struct {
	fileSet                *token.FileSet
	scannedGoFiles         []string
	violations             []violation
	interfaces             map[string]interfaceDecl
	mocks                  map[string][]methodDecl
	implementations        map[string][]methodDecl
	implementationBindings []implementationBinding
}

var placeholderReturnNamePattern = regexp.MustCompile(`^result[0-9]+$`)

var directDomainIdentifierConstructors = map[string]string{
	"Username":       "ParseUsername",
	"ConversationID": "ParseConversationID or ConversationIDFromUsername",
	"IdentityID":     "ParseIdentityID",
}

/* -------------------------------------------- Main -------------------------------------------- */

func main() {
	directories := parseDirectoriesOrExit()
	state := newAnalysisState()

	for _, directory := range directories {
		state.walkDirectory(directory)
	}

	state.addCrossFileViolations(directories)
	state.violations = dedupeViolations(state.violations)
	sortViolations(state.violations)
	printViolationsAndExit(state.violations)
}

func parseDirectoriesOrExit() (directories []string) {
	if len(os.Args) < minRequiredArgs {
		fmt.Fprintln(os.Stderr, "usage: stylecheck <dir>...")
		os.Exit(usageExitCode)
	}

	return os.Args[1:]
}

func newAnalysisState() (state *analysisState) {
	return &analysisState{
		fileSet:                token.NewFileSet(),
		scannedGoFiles:         make([]string, 0),
		interfaces:             make(map[string]interfaceDecl),
		mocks:                  make(map[string][]methodDecl),
		implementations:        make(map[string][]methodDecl),
		implementationBindings: make([]implementationBinding, 0),
	}
}

func (state *analysisState) walkDirectory(directory string) {
	walkError := filepath.WalkDir(
		directory,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if shouldSkipDirectory(entry) {
				return filepath.SkipDir
			}

			if entry.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}

			state.processFile(path)
			return nil
		},
	)
	if walkError != nil {
		fmt.Fprintf(os.Stderr, "error walking %s: %v\n", directory, walkError)
	}
}

func shouldSkipDirectory(entry os.DirEntry) (found bool) {
	if !entry.IsDir() {
		return false
	}

	switch entry.Name() {
	case "vendor", ".git", "testdata":
		return true
	default:
		return false
	}
}

func (state *analysisState) processFile(path string) {
	file, parseError := parser.ParseFile(state.fileSet, path, nil, parser.ParseComments)
	if parseError != nil {
		fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", path, parseError)
		return
	}

	normalisedPath := normalisePath(path)
	state.scannedGoFiles = append(state.scannedGoFiles, normalisedPath)
	isTestFile := strings.HasSuffix(path, "_test.go")
	state.addPerFileViolations(file, normalisedPath, isTestFile)
}

func (state *analysisState) addPerFileViolations(
	file *ast.File,
	normalisedPath string,
	isTestFile bool,
) {
	state.violations = append(state.violations, checkNamedReturns(state.fileSet, file)...)
	state.violations = append(state.violations, checkNakedReturns(state.fileSet, file)...)
	state.violations = append(state.violations, checkTypeElision(state.fileSet, file)...)
	state.violations = append(
		state.violations,
		checkGoErrorHandlingStyle(state.fileSet, file, normalisedPath, isTestFile)...,
	)
	state.violations = append(
		state.violations,
		checkAdapterErrorWrapping(state.fileSet, file, normalisedPath, isTestFile)...,
	)
	state.violations = append(
		state.violations,
		checkInlineCommentStyle(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(
		state.violations,
		checkDirectDomainIdentifierCasts(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(state.violations, checkParamOrder(state.fileSet, file)...)
	state.violations = append(state.violations, checkConstructorOrder(state.fileSet, file)...)
	if !isTestFile {
		state.violations = append(state.violations, checkFileStructureOrder(state.fileSet, file)...)
	}

	state.violations = append(
		state.violations,
		checkServiceTypeNaming(state.fileSet, file, normalisedPath)...,
	)
	state.violations = append(
		state.violations,
		checkCRUDLOrder(state.fileSet, file, normalisedPath)...,
	)

	collectInterfaces(state.fileSet, file, normalisedPath, state.interfaces)
	collectMockMethods(state.fileSet, file, normalisedPath, state.mocks)
	collectImplementationMethods(state.fileSet, file, normalisedPath, state.implementations)
	collectImplementationBindings(
		state.fileSet,
		file,
		normalisedPath,
		&state.implementationBindings,
	)

	// Single-letter checks skip test files to reduce noise in table-driven
	// structures and assertion helpers.
	if !isTestFile {
		state.violations = append(state.violations, checkSingleLetterVars(state.fileSet, file)...)
	}
}

func (state *analysisState) addCrossFileViolations(scanRoots []string) {
	typeAwareViolations, typeAwareRan := collectTypeAwareDomainIdentifierCastViolations(
		scanRoots,
		state.scannedGoFiles,
	)
	if typeAwareRan {
		state.violations = append(state.violations, typeAwareViolations...)
	}

	state.violations = append(
		state.violations,
		checkMockOrderAgainstInterfaces(state.interfaces, state.mocks)...,
	)
	state.violations = append(
		state.violations,
		checkImplementationOrderAgainstInterfaces(
			state.interfaces,
			state.implementations,
			state.implementationBindings,
		)...,
	)
}

func sortViolations(violations []violation) {
	sort.Slice(violations, func(i int, j int) bool {
		if violations[i].position.Filename == violations[j].position.Filename {
			return violations[i].position.Line < violations[j].position.Line
		}
		return violations[i].position.Filename < violations[j].position.Filename
	})
}

func dedupeViolations(violations []violation) (deduped []violation) {
	seen := make(map[string]bool)
	deduped = make([]violation, 0, len(violations))

	for _, current := range violations {
		key := fmt.Sprintf(
			"%s:%d:%d|%s|%s",
			current.position.Filename,
			current.position.Line,
			current.position.Column,
			current.rule,
			current.message,
		)

		if seen[key] {
			continue
		}

		seen[key] = true
		deduped = append(deduped, current)
	}

	return deduped
}

func printViolationsAndExit(violations []violation) {
	if len(violations) == 0 {
		return
	}

	for _, current := range violations {
		fmt.Fprintf(os.Stderr, "%s: [%s] %s\n", current.position, current.rule, current.message)
	}
	os.Exit(1)
}

/* ------------------------------------------- Checks ------------------------------------------- */

// checkNamedReturns ensures all functions, methods, and interface methods
// use named return values (2.2).
func checkNamedReturns(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl:
			violations = append(
				violations,
				checkFuncReturns(fileSet, declaration.Name.Name, declaration.Type)...,
			)

		case *ast.InterfaceType:
			if declaration.Methods == nil {
				return true
			}
			for _, method := range declaration.Methods.List {
				funcType, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue // embedded interface, skip
				}
				methodName := "(anonymous)"
				if len(method.Names) > 0 {
					methodName = method.Names[0].Name
				}
				violations = append(violations, checkFuncReturns(fileSet, methodName, funcType)...)
			}
		}
		return true
	})
	return violations
}

// checkFuncReturns reports a violation if any return value is unnamed.
func checkFuncReturns(
	fileSet *token.FileSet,
	funcName string,
	funcType *ast.FuncType,
) (violations []violation) {
	if funcType.Results == nil || len(funcType.Results.List) == 0 {
		return nil
	}
	for _, field := range funcType.Results.List {
		if len(field.Names) == 0 {
			violations = append(violations, violation{
				position: fileSet.Position(funcType.Results.Pos()),
				rule:     "2.2",
				message:  fmt.Sprintf("function %q has unnamed return values", funcName),
			})
			return violations // one violation per function is enough
		}

		for _, name := range field.Names {
			if placeholderReturnNamePattern.MatchString(name.Name) {
				violations = append(violations, violation{
					position: fileSet.Position(name.Pos()),
					rule:     "2.2",
					message: fmt.Sprintf(
						"function %q uses placeholder return name %q",
						funcName,
						name.Name,
					),
				})
			}
		}
	}
	return violations
}

// checkTypeElision ensures each parameter has its own type declaration (2.2).
func checkTypeElision(fileSet *token.FileSet, file *ast.File) (violations []violation) {
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
				violations = append(violations, violation{
					position: fileSet.Position(field.Pos()),
					rule:     "2.2",
					message: fmt.Sprintf(
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

// checkDirectDomainIdentifierCasts enforces parser/constructor usage for key domain IDs (2.2).
func checkDirectDomainIdentifierCasts(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if strings.Contains(path, domainPackagePathSegment) {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok || len(callExpression.Args) != 1 {
			return true
		}

		selector, ok := callExpression.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		packageIdentifier, ok := selector.X.(*ast.Ident)
		if !ok || packageIdentifier.Name != "domain" {
			return true
		}

		recommendedConstructor, found := directDomainIdentifierConstructors[selector.Sel.Name]
		if !found {
			return true
		}

		violations = append(violations, violation{
			position: fileSet.Position(callExpression.Pos()),
			rule:     "2.2",
			message: fmt.Sprintf(
				"direct cast to domain.%s is disallowed; use %s",
				selector.Sel.Name,
				recommendedConstructor,
			),
		})
		return true
	})

	return violations
}

func collectTypeAwareDomainIdentifierCastViolations(
	rootDirectories []string,
	filePaths []string,
) (violations []violation, ran bool) {
	if len(filePaths) == 0 || len(rootDirectories) == 0 {
		return nil, false
	}

	requestedFilePaths := make(map[string]bool, len(filePaths))
	for _, filePath := range filePaths {
		normalisedPath := normalisePath(filePath)
		requestedFilePaths[normalisedPath] = true
	}

	for _, rootDirectory := range rootDirectories {
		normalisedRoot := normalisePath(rootDirectory)

		packageConfig := &packages.Config{
			Mode: packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles |
				packages.NeedSyntax |
				packages.NeedTypes |
				packages.NeedTypesInfo,
			Dir:   normalisedRoot,
			Tests: true,
		}

		loadedPackages, err := packages.Load(packageConfig, "./...")
		if err != nil || len(loadedPackages) == 0 {
			continue
		}

		ran = true

		for _, loadedPackage := range loadedPackages {
			if loadedPackage == nil ||
				loadedPackage.TypesInfo == nil ||
				loadedPackage.Fset == nil {
				continue
			}

			for _, file := range loadedPackage.Syntax {
				filePath := normalisePath(loadedPackage.Fset.Position(file.Pos()).Filename)

				if !requestedFilePaths[filePath] {
					continue
				}

				if strings.Contains(filePath, domainPackagePathSegment) {
					continue
				}

				ast.Inspect(file, func(node ast.Node) bool {
					callExpression, ok := node.(*ast.CallExpr)
					if !ok || len(callExpression.Args) != 1 {
						return true
					}

					typeAndValue, ok := loadedPackage.TypesInfo.Types[callExpression.Fun]
					if !ok {
						return true
					}

					domainTypeName, found := resolvedDomainIdentifierTypeName(typeAndValue.Type)
					if !found {
						return true
					}

					recommendedConstructor := directDomainIdentifierConstructors[domainTypeName]
					violations = append(violations, violation{
						position: loadedPackage.Fset.Position(callExpression.Pos()),
						rule:     "2.2",
						message: fmt.Sprintf(
							"direct cast to domain.%s is disallowed; use %s",
							domainTypeName,
							recommendedConstructor,
						),
					})

					return true
				})
			}
		}
	}

	return violations, ran
}

func resolvedDomainIdentifierTypeName(targetType types.Type) (name string, found bool) {
	namedType, ok := types.Unalias(targetType).(*types.Named)
	if !ok {
		return "", false
	}

	typeObject := namedType.Obj()
	if typeObject == nil || typeObject.Pkg() == nil {
		return "", false
	}

	packagePath := typeObject.Pkg().Path()
	if !isDomainPackagePath(packagePath) {
		return "", false
	}

	typeName := typeObject.Name()
	if _, supported := directDomainIdentifierConstructors[typeName]; !supported {
		return "", false
	}

	return typeName, true
}

func isDomainPackagePath(packagePath string) (found bool) {
	if packagePath == "internal/core/domain" {
		return true
	}

	return strings.HasSuffix(packagePath, domainPackagePathSuffix)
}

func normalisePath(path string) (normalisedPath string) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(filepath.Clean(absolutePath))
}

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

// checkNakedReturns reports naked returns in functions that declare named return values (2.2).
func checkNakedReturns(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		functionDecl, ok := node.(*ast.FuncDecl)
		if !ok ||
			functionDecl.Type == nil ||
			functionDecl.Type.Results == nil ||
			functionDecl.Body == nil {
			return true
		}

		if !funcHasNamedReturns(functionDecl.Type) {
			return true
		}

		ast.Inspect(functionDecl.Body, func(bodyNode ast.Node) bool {
			switch typed := bodyNode.(type) {
			case *ast.FuncLit:
				return false

			case *ast.ReturnStmt:
				if len(typed.Results) == 0 {
					violations = append(violations, violation{
						position: fileSet.Position(typed.Pos()),
						rule:     "2.2",
						message: fmt.Sprintf(
							"function %q uses a naked return",
							functionDecl.Name.Name,
						),
					})
				}
			}

			return true
		})

		return true
	})

	return violations
}

func funcHasNamedReturns(functionType *ast.FuncType) (found bool) {
	for _, resultField := range functionType.Results.List {
		if len(resultField.Names) > 0 {
			return true
		}
	}

	return false
}

// checkSingleLetterVars flags single-letter variable names that are not
// loop indices (i, j, k) or method receivers (2.2).
func checkSingleLetterVars(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	allowed := map[string]bool{"i": true, "j": true, "k": true, "_": true}

	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl:
			violations = append(
				violations,
				checkSingleLetterFuncParams(fileSet, declaration, allowed)...,
			)

		case *ast.AssignStmt:
			violations = append(
				violations,
				checkSingleLetterAssignStmt(fileSet, declaration, allowed)...,
			)

		case *ast.RangeStmt:
			violations = append(
				violations,
				checkSingleLetterRangeStmt(fileSet, declaration, allowed)...,
			)

		case *ast.ValueSpec:
			violations = append(
				violations,
				checkSingleLetterValueSpec(fileSet, declaration, allowed)...,
			)
		}
		return true
	})

	return violations
}

func checkSingleLetterFuncParams(
	fileSet *token.FileSet,
	declaration *ast.FuncDecl,
	allowed map[string]bool,
) (violations []violation) {
	// Check parameters only (Recv is excluded, so receivers pass).
	if declaration.Type.Params == nil {
		return nil
	}

	for _, field := range declaration.Type.Params.List {
		for _, name := range field.Names {
			violationValue := singleLetterNameViolation(
				fileSet,
				name,
				allowed,
				"2.2",
				fmt.Sprintf(
					"single-letter parameter %q in function %q",
					name.Name,
					declaration.Name.Name,
				),
			)
			if violationValue != nil {
				violations = append(violations, *violationValue)
			}
		}
	}

	return violations
}

func checkSingleLetterAssignStmt(
	fileSet *token.FileSet,
	declaration *ast.AssignStmt,
	allowed map[string]bool,
) (violations []violation) {
	if declaration.Tok != token.DEFINE {
		return nil
	}

	for _, lhs := range declaration.Lhs {
		identifierNode, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}

		violationValue := singleLetterNameViolation(
			fileSet,
			identifierNode,
			allowed,
			"2.2",
			fmt.Sprintf("single-letter variable %q", identifierNode.Name),
		)
		if violationValue != nil {
			violations = append(violations, *violationValue)
		}
	}

	return violations
}

func checkSingleLetterRangeStmt(
	fileSet *token.FileSet,
	declaration *ast.RangeStmt,
	allowed map[string]bool,
) (violations []violation) {
	if declaration.Tok != token.DEFINE {
		return nil
	}

	if key, ok := declaration.Key.(*ast.Ident); ok {
		violationValue := singleLetterNameViolation(
			fileSet,
			key,
			allowed,
			"2.2",
			fmt.Sprintf("single-letter range variable %q", key.Name),
		)
		if violationValue != nil {
			violations = append(violations, *violationValue)
		}
	}

	if declaration.Value != nil {
		if value, ok := declaration.Value.(*ast.Ident); ok {
			violationValue := singleLetterNameViolation(
				fileSet,
				value,
				allowed,
				"2.2",
				fmt.Sprintf("single-letter range variable %q", value.Name),
			)
			if violationValue != nil {
				violations = append(violations, *violationValue)
			}
		}
	}

	return violations
}

func checkSingleLetterValueSpec(
	fileSet *token.FileSet,
	declaration *ast.ValueSpec,
	allowed map[string]bool,
) (violations []violation) {
	for _, name := range declaration.Names {
		violationValue := singleLetterNameViolation(
			fileSet,
			name,
			allowed,
			"2.2",
			fmt.Sprintf("single-letter variable %q", name.Name),
		)
		if violationValue != nil {
			violations = append(violations, *violationValue)
		}
	}

	return violations
}

func singleLetterNameViolation(
	fileSet *token.FileSet,
	name *ast.Ident,
	allowed map[string]bool,
	rule string,
	message string,
) (violationValue *violation) {
	if len(name.Name) != 1 || allowed[name.Name] {
		return nil
	}

	return &violation{
		position: fileSet.Position(name.Pos()),
		rule:     rule,
		message:  message,
	}
}

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

		// ctx must be first if present.
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

		// Secrets must be last.
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

// checkConstructorOrder ensures constructor parameters follow the canonical
// ordering: repositories -> services -> adapters -> config -> secrets (2.8).
func checkConstructorOrder(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	categoryLabel := map[int]string{
		categoryRepository: "repository",
		categoryService:    "service",
		categoryAdapter:    "adapter",
		categoryConfig:     "config",
		categorySecret:     "secret",
	}

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
						categoryLabel[category],
						categoryLabel[prevCategory],
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

// checkFileStructureOrder enforces objective top-level declaration ordering (2.9).
// This check intentionally avoids subjective formatting requirements.
func checkFileStructureOrder(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	highestSeenCategory := declUnknown
	seenCategories := map[int]bool{}

	for _, declaration := range file.Decls {
		currentCategory := classifyTopLevelDecl(declaration)
		if currentCategory == declUnknown {
			continue
		}

		seenCategories[currentCategory] = true
		if currentCategory < highestSeenCategory {
			violations = append(violations, violation{
				position: fileSet.Position(declaration.Pos()),
				rule:     "2.9",
				message: fmt.Sprintf(
					"declaration group %q appears after %q",
					declCategoryName(currentCategory),
					declCategoryName(highestSeenCategory),
				),
			})
			continue
		}

		highestSeenCategory = currentCategory
	}

	if len(seenCategories) < minCategoryKinds {
		return nil
	}
	return violations
}

// checkServiceTypeNaming enforces exported type naming in internal/core/services (2.2).
func checkServiceTypeNaming(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if !strings.Contains(path, "/internal/core/services/") {
		return nil
	}
	// accountref is a parsing helper package, not a use-case/service package.
	if strings.Contains(path, "/internal/core/services/accountref/") {
		return nil
	}

	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || !typeSpec.Name.IsExported() {
				continue
			}
			name := typeSpec.Name.Name
			if strings.HasSuffix(name, "Service") ||
				strings.HasSuffix(name, "UseCase") ||
				strings.HasSuffix(name, "Config") {
				continue
			}
			violations = append(violations, violation{
				position: fileSet.Position(typeSpec.Pos()),
				rule:     "2.2",
				message: fmt.Sprintf(
					"exported type %q should end with Service, UseCase, or Config",
					name,
				),
			})
		}
	}
	return violations
}

// checkCRUDLOrder validates CRUD-L method ordering inside ports interfaces (2.5).
func checkCRUDLOrder(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if !strings.Contains(path, "/internal/core/ports/") {
		return nil
	}

	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok || interfaceType.Methods == nil {
				continue
			}

			lastCategory := crudUnknown
			lastMethod := ""

			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue // embedded interface
				}
				name := method.Names[0].Name
				category := crudCategory(name)
				if category == crudUnknown {
					continue
				}
				if lastCategory != crudUnknown && category < lastCategory {
					violations = append(violations, violation{
						position: fileSet.Position(method.Pos()),
						rule:     "2.5",
						message: fmt.Sprintf(
							"method %q in interface %q is out of CRUD-L order (after %q)",
							name,
							typeSpec.Name.Name,
							lastMethod,
						),
					})
				}
				lastCategory = category
				lastMethod = name
			}
		}
	}
	return violations
}

// collectInterfaces records interface method order from internal/core/ports.
func collectInterfaces(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	interfaces map[string]interfaceDecl,
) {
	if !strings.Contains(path, "/internal/core/ports/") {
		return
	}

	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok || interfaceType.Methods == nil {
				continue
			}

			methods := make([]methodDecl, 0, len(interfaceType.Methods.List))
			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue // embedded interface
				}
				methods = append(methods, methodDecl{
					name:     method.Names[0].Name,
					position: fileSet.Position(method.Pos()),
				})
			}

			interfaces[typeSpec.Name.Name] = interfaceDecl{
				name:     typeSpec.Name.Name,
				methods:  methods,
				position: fileSet.Position(typeSpec.Pos()),
			}
		}
	}
}

// collectMockMethods records method order for each mock receiver type.
func collectMockMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	mocks map[string][]methodDecl,
) {
	if !strings.Contains(path, "/internal/mocks/") {
		return
	}
	for _, declaration := range file.Decls {
		funcDecl, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}
		receiver := receiverTypeName(funcDecl.Recv.List[0].Type)
		if receiver == "" {
			continue
		}
		mocks[receiver] = append(mocks[receiver], methodDecl{
			name:     funcDecl.Name.Name,
			position: fileSet.Position(funcDecl.Name.Pos()),
		})
	}
}

// checkMockOrderAgainstInterfaces compares mock method order with ports interface order (2.5).
func checkMockOrderAgainstInterfaces(
	interfaces map[string]interfaceDecl,
	mocks map[string][]methodDecl,
) (violations []violation) {
	interfaceNames := make([]string, 0, len(interfaces))
	for name := range interfaces {
		interfaceNames = append(interfaceNames, name)
	}
	sort.Strings(interfaceNames)

	for _, interfaceName := range interfaceNames {
		interfaceDecl := interfaces[interfaceName]
		mockMethods, matchedMockName, ambiguousMockNames, found := resolveMockMethodsForInterface(
			interfaceName,
			mocks,
		)
		if len(ambiguousMockNames) > 0 {
			position := interfaceDecl.position
			if len(mockMethods) > 0 {
				position = mockMethods[0].position
			}
			violations = append(violations, violation{
				position: position,
				rule:     "2.5",
				message: fmt.Sprintf(
					"multiple mock types match interface %q: %s",
					interfaceName,
					strings.Join(ambiguousMockNames, ", "),
				),
			})
			continue
		}

		if !found {
			continue // only compare where a mock exists
		}

		interfaceMethodNames := make([]string, len(interfaceDecl.methods))
		for i, method := range interfaceDecl.methods {
			interfaceMethodNames[i] = method.name
		}
		mockMethodNames := make([]string, len(mockMethods))
		for i, method := range mockMethods {
			mockMethodNames[i] = method.name
		}

		if len(interfaceMethodNames) != len(mockMethodNames) {
			pos := interfaceDecl.position
			if len(mockMethods) > 0 {
				pos = mockMethods[0].position
			}
			violations = append(violations, violation{
				position: pos,
				rule:     "2.5",
				message: fmt.Sprintf(
					"mock %q for interface %q method count (%d) does not match interface (%d)",
					matchedMockName,
					interfaceName,
					len(mockMethodNames),
					len(interfaceMethodNames),
				),
			})
			continue
		}

		for index := range interfaceMethodNames {
			if interfaceMethodNames[index] == mockMethodNames[index] {
				continue
			}
			violations = append(violations, violation{
				position: mockMethods[index].position,
				rule:     "2.5",
				message: fmt.Sprintf(
					"mock %q for interface %q method order mismatch at position %d: "+
						"got %q, want %q",
					matchedMockName,
					interfaceName,
					index+1,
					mockMethodNames[index],
					interfaceMethodNames[index],
				),
			})
			break
		}
	}

	return violations
}

// checkImplementationOrderAgainstInterfaces compares implementation method order with
// ports interface order for types that declare compile-time assertions (2.5).
func checkImplementationOrderAgainstInterfaces(
	interfaces map[string]interfaceDecl,
	implementations map[string][]methodDecl,
	bindings []implementationBinding,
) (violations []violation) {
	sort.Slice(bindings, func(i int, j int) bool {
		if bindings[i].interfaceName == bindings[j].interfaceName {
			return bindings[i].implementationName < bindings[j].implementationName
		}
		return bindings[i].interfaceName < bindings[j].interfaceName
	})

	for _, binding := range bindings {
		interfaceDeclaration, found := interfaces[binding.interfaceName]
		if !found {
			continue
		}

		implementationMethods, found := implementations[binding.implementationKey]
		if !found {
			continue
		}

		interfaceMethodNames := make([]string, len(interfaceDeclaration.methods))
		interfaceMethodNamesSet := make(map[string]bool, len(interfaceDeclaration.methods))
		for i, method := range interfaceDeclaration.methods {
			interfaceMethodNames[i] = method.name
			interfaceMethodNamesSet[method.name] = true
		}

		implementationInterfaceMethods := make([]methodDecl, 0, len(interfaceMethodNames))
		for _, method := range implementationMethods {
			if interfaceMethodNamesSet[method.name] {
				implementationInterfaceMethods = append(implementationInterfaceMethods, method)
			}
		}

		if len(implementationInterfaceMethods) != len(interfaceMethodNames) {
			violations = append(violations, violation{
				position: binding.position,
				rule:     "2.5",
				message: fmt.Sprintf(
					"implementation %q for interface %q method count (%d) "+
						"does not match interface (%d)",
					binding.implementationName,
					binding.interfaceName,
					len(implementationInterfaceMethods),
					len(interfaceMethodNames),
				),
			})
			continue
		}

		for index := range interfaceMethodNames {
			if implementationInterfaceMethods[index].name == interfaceMethodNames[index] {
				continue
			}

			violations = append(violations, violation{
				position: implementationInterfaceMethods[index].position,
				rule:     "2.5",
				message: fmt.Sprintf(
					"implementation %q for interface %q method order mismatch at position %d: "+
						"got %q, want %q",
					binding.implementationName,
					binding.interfaceName,
					index+1,
					implementationInterfaceMethods[index].name,
					interfaceMethodNames[index],
				),
			})
			break
		}
	}

	return violations
}

/* ------------------------------------- Resolution Helpers ------------------------------------- */

func resolveMockMethodsForInterface(
	interfaceName string,
	mocks map[string][]methodDecl,
) (
	methods []methodDecl,
	matchedMockName string,
	ambiguousMockNames []string,
	found bool,
) {
	if directMethods, ok := mocks[interfaceName]; ok {
		return directMethods, interfaceName, nil, true
	}

	interfaceCanonicalName := normaliseMockTypeName(interfaceName)
	for mockTypeName, mockMethods := range mocks {
		if normaliseMockTypeName(mockTypeName) != interfaceCanonicalName {
			continue
		}

		ambiguousMockNames = append(ambiguousMockNames, mockTypeName)
		if len(methods) == 0 {
			methods = mockMethods
			matchedMockName = mockTypeName
		}
	}

	if len(ambiguousMockNames) == 0 {
		return nil, "", nil, false
	}

	sort.Strings(ambiguousMockNames)
	if len(ambiguousMockNames) > 1 {
		return methods, matchedMockName, ambiguousMockNames, false
	}

	return methods, matchedMockName, nil, true
}

func collectImplementationMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	implementations map[string][]methodDecl,
) {
	isPortPath := strings.Contains(path, "/internal/core/ports/")
	isMockPath := strings.Contains(path, "/internal/mocks/")
	if isPortPath || isMockPath {
		return
	}

	for _, declaration := range file.Decls {
		funcDeclaration, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDeclaration.Recv == nil || len(funcDeclaration.Recv.List) == 0 {
			continue
		}

		receiverName := receiverTypeName(funcDeclaration.Recv.List[0].Type)
		if receiverName == "" {
			continue
		}

		key := typeDeclKey(path, receiverName)
		implementations[key] = append(implementations[key], methodDecl{
			name:     funcDeclaration.Name.Name,
			position: fileSet.Position(funcDeclaration.Name.Pos()),
		})
	}
}

func collectImplementationBindings(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	bindings *[]implementationBinding,
) {
	if strings.Contains(path, "/internal/mocks/") {
		return
	}

	for _, declaration := range file.Decls {
		genDeclaration, ok := declaration.(*ast.GenDecl)
		if !ok || genDeclaration.Tok != token.VAR {
			continue
		}

		for _, spec := range genDeclaration.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			if len(valueSpec.Names) != 1 || valueSpec.Names[0].Name != "_" {
				continue
			}

			interfaceName := typeNameFromExpr(valueSpec.Type)
			if interfaceName == "" || len(valueSpec.Values) != 1 {
				continue
			}

			implementationName := implementationTypeFromAssertion(valueSpec.Values[0])
			if implementationName == "" {
				continue
			}

			*bindings = append(*bindings, implementationBinding{
				interfaceName:      interfaceName,
				implementationName: implementationName,
				implementationKey:  typeDeclKey(path, implementationName),
				position:           fileSet.Position(valueSpec.Pos()),
			})
		}
	}
}

/* -------------------------------------- AST Type Helpers -------------------------------------- */

func typeNameFromExpr(expression ast.Expr) (name string) {
	switch typed := expression.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.SelectorExpr:
		return typed.Sel.Name
	default:
		return ""
	}
}

func implementationTypeFromAssertion(expression ast.Expr) (name string) {
	switch typed := expression.(type) {
	case *ast.CallExpr:
		return implementationTypeFromAssertion(typed.Fun)

	case *ast.ParenExpr:
		return implementationTypeFromAssertion(typed.X)

	case *ast.StarExpr:
		return typeNameFromExpr(typed.X)

	case *ast.UnaryExpr:
		if typed.Op == token.AND {
			return implementationTypeFromAssertion(typed.X)
		}
		return ""

	case *ast.CompositeLit:
		return typeNameFromExpr(typed.Type)

	case *ast.Ident:
		return typed.Name

	default:
		return ""
	}
}

func typeDeclKey(path string, typeName string) (key string) {
	return fmt.Sprintf("%s::%s", filepath.ToSlash(filepath.Dir(path)), typeName)
}

func normaliseMockTypeName(typeName string) (normalisedTypeName string) {
	normalisedTypeName = typeName

	for strings.HasPrefix(normalisedTypeName, "Mock") {
		normalisedTypeName = strings.TrimPrefix(normalisedTypeName, "Mock")
	}

	for strings.HasSuffix(normalisedTypeName, "Mock") {
		normalisedTypeName = strings.TrimSuffix(normalisedTypeName, "Mock")
	}

	return normalisedTypeName
}

/* ----------------------------------- Classification Helpers ----------------------------------- */

// isSecretName returns true if the parameter name represents a secret.
func isSecretName(name string) (found bool) {
	secrets := map[string]bool{
		"passphrase": true, "privateKey": true, "token": true,
		"seed": true, "secret": true, "password": true, "secretKey": true,
	}
	return secrets[name]
}

// isConfigName returns true if the parameter name represents configuration.
func isConfigName(name string) (found bool) {
	configs := map[string]bool{
		"serverURL": true, "relayURL": true, "identityID": true, "timeout": true,
	}
	return configs[name]
}

// isConstructor returns true if the function name follows the NewXxx pattern.
func isConstructor(name string) (found bool) {
	if name == "New" {
		return true
	}
	return strings.HasPrefix(name, "New") && len(name) > 3 && unicode.IsUpper(rune(name[3]))
}

func classifyTopLevelDecl(declaration ast.Decl) (category int) {
	switch typed := declaration.(type) {
	case *ast.GenDecl:
		switch typed.Tok {
		case token.CONST:
			return declConstants
		case token.TYPE:
			return declTypes
		case token.VAR:
			if isCompileTimeAssertionDecl(typed) {
				return declAssertions
			}

			if isSentinelErrorDecl(typed) {
				return declErrors
			}
		}
	}

	return declUnknown
}

func isCompileTimeAssertionDecl(declaration *ast.GenDecl) (found bool) {
	if declaration.Tok != token.VAR || len(declaration.Specs) == 0 {
		return false
	}

	for _, spec := range declaration.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok || len(valueSpec.Names) != 1 || valueSpec.Names[0].Name != "_" {
			return false
		}

		if valueSpec.Type == nil || len(valueSpec.Values) != 1 {
			return false
		}
	}

	return true
}

func isSentinelErrorDecl(declaration *ast.GenDecl) (found bool) {
	if declaration.Tok != token.VAR || len(declaration.Specs) == 0 {
		return false
	}

	for _, spec := range declaration.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok || len(valueSpec.Names) == 0 {
			return false
		}
		for _, name := range valueSpec.Names {
			if !strings.HasPrefix(name.Name, "Err") {
				return false
			}
		}
	}

	return true
}

func declCategoryName(category int) (name string) {
	switch category {
	case declConstants:
		return "constants"
	case declErrors:
		return "errors"
	case declTypes:
		return "types"
	case declAssertions:
		return "assertions"
	default:
		return "unknown"
	}
}

// classifyParam determines the category of a constructor parameter.
func classifyParam(field *ast.Field) (category int) {
	typeName := typeString(field.Type)

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

// crudCategory classifies method names into CRUD-L categories.
func crudCategory(name string) (category int) {
	switch {
	case strings.HasPrefix(name, "List"):
		return crudList
	case strings.HasPrefix(name, "Delete"),
		strings.HasPrefix(name, "Remove"),
		strings.HasPrefix(name, "Consume"):
		return crudDelete
	case strings.HasPrefix(name, "Update"),
		strings.HasPrefix(name, "Set"),
		strings.HasPrefix(name, "Ack"):
		return crudUpdate
	case strings.HasPrefix(name, "Read"),
		strings.HasPrefix(name, "Load"),
		strings.HasPrefix(name, "Get"),
		strings.HasPrefix(name, "Fetch"),
		strings.HasPrefix(name, "IdentityExists"),
		strings.HasPrefix(name, "Metadata"),
		strings.HasPrefix(name, "Fingerprint"),
		strings.HasPrefix(name, "Current"):
		return crudRead
	case strings.HasPrefix(name, "Create"),
		strings.HasPrefix(name, "Save"),
		strings.HasPrefix(name, "Generate"),
		strings.HasPrefix(name, "Register"),
		strings.HasPrefix(name, "Initiate"),
		strings.HasPrefix(name, "Send"):
		return crudCreate
	default:
		return crudUnknown
	}
}

// receiverTypeName returns the receiver type for methods (supports T and *T).
func receiverTypeName(expr ast.Expr) (typeName string) {
	switch typed := expr.(type) {
	case *ast.StarExpr:
		if identifierNode, ok := typed.X.(*ast.Ident); ok {
			return identifierNode.Name
		}
	case *ast.Ident:
		return typed.Name
	}
	return ""
}

// typeString extracts a human-readable type name from an AST expression.
func typeString(expression ast.Expr) (typeName string) {
	switch typed := expression.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.SelectorExpr:
		return typeString(typed.X) + "." + typed.Sel.Name
	case *ast.StarExpr:
		return typeString(typed.X)
	case *ast.ArrayType:
		return "[]" + typeString(typed.Elt)
	default:
		return ""
	}
}
