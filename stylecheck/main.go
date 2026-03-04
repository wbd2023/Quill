// Command stylecheck performs AST-based style checks on Go source files.
//
// It enforces the following rules from STYLE.md:
//   - 2.2 Named returns: all functions must use named return values.
//   - 2.2 Type elision: each parameter must have its own type.
//   - 2.2 Single-letter variables: only i, j, k (loops) and receivers.
//   - 2.2 Service package type naming: exported types end with Service/UseCase/Config.
//   - 2.5 CRUD-L ordering inside interfaces.
//   - 2.5 Mock method order matches interface method order exactly.
//   - 2.7 Parameter ordering: ctx first, secrets last.
//   - 2.8 Constructor ordering: repos -> services -> adapters -> config -> secrets.
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
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

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
	minRequiredArgs   = 2
	usageExitCode     = 2
	minParamFieldSpan = 2
)

/* -------------------------------------------- Main -------------------------------------------- */

func main() {
	if len(os.Args) < minRequiredArgs {
		fmt.Fprintln(os.Stderr, "usage: stylecheck <dir>...")
		os.Exit(usageExitCode)
	}

	directories := os.Args[1:]
	var allViolations []violation
	fileSet := token.NewFileSet()
	interfaces := make(map[string]interfaceDecl)
	mocks := make(map[string][]methodDecl)

	for _, directory := range directories {
		walkError := filepath.WalkDir(directory, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if entry.IsDir() {
				switch entry.Name() {
				case "vendor", ".git", "testdata":
					return filepath.SkipDir
				}
				return nil
			}

			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			file, parseError := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
			if parseError != nil {
				fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", path, parseError)
				return nil
			}

			normalisedPath := filepath.ToSlash(path)
			isTestFile := strings.HasSuffix(path, "_test.go")

			allViolations = append(allViolations, checkNamedReturns(fileSet, file)...)
			allViolations = append(allViolations, checkTypeElision(fileSet, file)...)
			allViolations = append(allViolations, checkParamOrder(fileSet, file)...)
			allViolations = append(allViolations, checkConstructorOrder(fileSet, file)...)
			allViolations = append(
				allViolations,
				checkServiceTypeNaming(fileSet, file, normalisedPath)...,
			)
			allViolations = append(
				allViolations,
				checkCRUDLOrder(fileSet, file, normalisedPath)...,
			)
			collectInterfaces(fileSet, file, normalisedPath, interfaces)
			collectMockMethods(fileSet, file, normalisedPath, mocks)

			// Single-letter checks skip test files to reduce noise from
			// table-driven test structures and assertion helpers.
			if !isTestFile {
				allViolations = append(allViolations, checkSingleLetterVars(fileSet, file)...)
			}

			return nil
		})
		if walkError != nil {
			fmt.Fprintf(os.Stderr, "error walking %s: %v\n", directory, walkError)
		}
	}

	allViolations = append(allViolations, checkMockOrderAgainstInterfaces(interfaces, mocks)...)
	sort.Slice(allViolations, func(i int, j int) bool {
		if allViolations[i].position.Filename == allViolations[j].position.Filename {
			return allViolations[i].position.Line < allViolations[j].position.Line
		}
		return allViolations[i].position.Filename < allViolations[j].position.Filename
	})

	if len(allViolations) > 0 {
		for _, current := range allViolations {
			fmt.Fprintf(os.Stderr, "%s: [%s] %s\n", current.position, current.rule, current.message)
		}
		os.Exit(1)
	}
}

/* ------------------------------------------- Checks ------------------------------------------- */

// checkNamedReturns ensures all functions, methods, and interface methods
// use named return values (2.2).
func checkNamedReturns(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl:
			violations = append(violations, checkFuncReturns(fileSet, declaration.Name.Name, declaration.Type)...)

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
func checkFuncReturns(fileSet *token.FileSet, funcName string, funcType *ast.FuncType) (violations []violation) {
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
	}
	return nil
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
					message:  fmt.Sprintf("type elision: parameters %s share a type", strings.Join(names, ", ")),
				})
			}
		}
		return true
	})
	return violations
}

// checkSingleLetterVars flags single-letter variable names that are not
// loop indices (i, j, k) or method receivers (2.2).
func checkSingleLetterVars(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	allowed := map[string]bool{"i": true, "j": true, "k": true, "_": true}

	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl:
			// Check parameters only (Recv is excluded, so receivers pass).
			if declaration.Type.Params != nil {
				for _, field := range declaration.Type.Params.List {
					for _, name := range field.Names {
						if len(name.Name) == 1 && !allowed[name.Name] {
							violations = append(violations, violation{
								position: fileSet.Position(name.Pos()),
								rule:     "2.2",
								message:  fmt.Sprintf("single-letter parameter %q in function %q", name.Name, declaration.Name.Name),
							})
						}
					}
				}
			}

		case *ast.AssignStmt:
			if declaration.Tok != token.DEFINE {
				return true
			}
			for _, lhs := range declaration.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}
				if len(ident.Name) == 1 && !allowed[ident.Name] {
					violations = append(violations, violation{
						position: fileSet.Position(ident.Pos()),
						rule:     "2.2",
						message:  fmt.Sprintf("single-letter variable %q", ident.Name),
					})
				}
			}

		case *ast.RangeStmt:
			if declaration.Tok != token.DEFINE {
				return true
			}

			if key, ok := declaration.Key.(*ast.Ident); ok {
				if len(key.Name) == 1 && !allowed[key.Name] {
					violations = append(violations, violation{
						position: fileSet.Position(key.Pos()),
						rule:     "2.2",
						message:  fmt.Sprintf("single-letter range variable %q", key.Name),
					})
				}
			}
			if declaration.Value != nil {
				if value, ok := declaration.Value.(*ast.Ident); ok {
					if len(value.Name) == 1 && !allowed[value.Name] {
						violations = append(violations, violation{
							position: fileSet.Position(value.Pos()),
							rule:     "2.2",
							message:  fmt.Sprintf("single-letter range variable %q", value.Name),
						})
					}
				}
			}

		case *ast.ValueSpec:
			for _, name := range declaration.Names {
				if len(name.Name) == 1 && !allowed[name.Name] {
					violations = append(violations, violation{
						position: fileSet.Position(name.Pos()),
						rule:     "2.2",
						message:  fmt.Sprintf("single-letter variable %q", name.Name),
					})
				}
			}
		}
		return true
	})
	return violations
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
					"mock %q for interface %q method order mismatch at position %d: got %q, want %q",
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

/* ------------------------------------------- Helpers ------------------------------------------ */

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
		if ident, ok := typed.X.(*ast.Ident); ok {
			return ident.Name
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
