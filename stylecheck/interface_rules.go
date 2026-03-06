package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
)

/* --------------------------------------- Interface Rules -------------------------------------- */

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
					continue
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
			continue
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
			position := interfaceDecl.position
			if len(mockMethods) > 0 {
				position = mockMethods[0].position
			}
			violations = append(violations, violation{
				position: position,
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
