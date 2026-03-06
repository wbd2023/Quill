package checker

import (
	"fmt"
	"sort"
	"strings"

	"stylecheck/internal/checker/collect"
)

/* --------------------------------------- Interface Rules -------------------------------------- */

// checkMockOrderAgainstInterfaces compares mock method order with ports interface order (2.5).
func checkMockOrderAgainstInterfaces(
	interfaces map[string]collect.InterfaceDecl,
	mocks map[string][]collect.MethodDecl,
) (violations []violation) {
	interfaceNames := make([]string, 0, len(interfaces))
	for name := range interfaces {
		interfaceNames = append(interfaceNames, name)
	}
	sort.Strings(interfaceNames)

	for _, interfaceName := range interfaceNames {
		interfaceDeclaration := interfaces[interfaceName]
		mockMethods, matchedMockName, ambiguousMockNames, found := resolveMockMethodsForInterface(
			interfaceName,
			mocks,
		)
		if len(ambiguousMockNames) > 0 {
			position := interfaceDeclaration.Position
			if len(mockMethods) > 0 {
				position = mockMethods[0].Position
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

		interfaceMethodNames := make([]string, len(interfaceDeclaration.Methods))
		for i, method := range interfaceDeclaration.Methods {
			interfaceMethodNames[i] = method.Name
		}
		mockMethodNames := make([]string, len(mockMethods))
		for i, method := range mockMethods {
			mockMethodNames[i] = method.Name
		}

		if len(interfaceMethodNames) != len(mockMethodNames) {
			position := interfaceDeclaration.Position
			if len(mockMethods) > 0 {
				position = mockMethods[0].Position
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
				position: mockMethods[index].Position,
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

// checkImplementationOrderAgainstInterfaces compares implementation method order with ports
// interface order for types that declare compile-time assertions (2.5).
func checkImplementationOrderAgainstInterfaces(
	interfaces map[string]collect.InterfaceDecl,
	implementations map[string][]collect.MethodDecl,
	bindings []collect.ImplementationBinding,
) (violations []violation) {
	sort.Slice(bindings, func(i int, j int) bool {
		if bindings[i].InterfaceName == bindings[j].InterfaceName {
			return bindings[i].ImplementationName < bindings[j].ImplementationName
		}
		return bindings[i].InterfaceName < bindings[j].InterfaceName
	})

	for _, binding := range bindings {
		interfaceDeclaration, found := interfaces[binding.InterfaceName]
		if !found {
			continue
		}

		implementationMethods, found := implementations[binding.ImplementationKey]
		if !found {
			continue
		}

		interfaceMethodNames := make([]string, len(interfaceDeclaration.Methods))
		interfaceMethodNamesSet := make(map[string]bool, len(interfaceDeclaration.Methods))
		for i, method := range interfaceDeclaration.Methods {
			interfaceMethodNames[i] = method.Name
			interfaceMethodNamesSet[method.Name] = true
		}

		implementationInterfaceMethods := make([]collect.MethodDecl, 0, len(interfaceMethodNames))
		for _, method := range implementationMethods {
			if interfaceMethodNamesSet[method.Name] {
				implementationInterfaceMethods = append(implementationInterfaceMethods, method)
			}
		}

		if len(implementationInterfaceMethods) != len(interfaceMethodNames) {
			violations = append(violations, violation{
				position: binding.Position,
				rule:     "2.5",
				message: fmt.Sprintf(
					"implementation %q for interface %q method count (%d) "+
						"does not match interface (%d)",
					binding.ImplementationName,
					binding.InterfaceName,
					len(implementationInterfaceMethods),
					len(interfaceMethodNames),
				),
			})
			continue
		}

		for index := range interfaceMethodNames {
			if implementationInterfaceMethods[index].Name == interfaceMethodNames[index] {
				continue
			}

			violations = append(violations, violation{
				position: implementationInterfaceMethods[index].Position,
				rule:     "2.5",
				message: fmt.Sprintf(
					"implementation %q for interface %q method order mismatch at position %d: "+
						"got %q, want %q",
					binding.ImplementationName,
					binding.InterfaceName,
					index+1,
					implementationInterfaceMethods[index].Name,
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
	mocks map[string][]collect.MethodDecl,
) (
	methods []collect.MethodDecl,
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
