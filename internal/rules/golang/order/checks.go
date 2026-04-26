package order

import (
	"fmt"
	"sort"
	"strings"

	"ciphera/tools/internal/rules/golang/checks"
)

/* --------------------------------------- Interface Rules -------------------------------------- */

func checkMockOrderAgainstInterfaces(
	interfaces map[string]interfaceDeclaration,
	mocks map[string][]methodDeclaration,
) (violations []checks.Violation) {
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
			violations = append(violations, checks.Violation{
				Position: position,
				Rule:     checks.DiagnosticMockOrder,
				Message: fmt.Sprintf(
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
			violations = append(violations, checks.Violation{
				Position: position,
				Rule:     checks.DiagnosticMockOrder,
				Message: fmt.Sprintf(
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
			violations = append(violations, checks.Violation{
				Position: mockMethods[index].Position,
				Rule:     checks.DiagnosticMockOrder,
				Message: fmt.Sprintf(
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

func checkImplementationOrderAgainstInterfaces(
	interfaces map[string]interfaceDeclaration,
	implementations map[string][]methodDeclaration,
	bindings []implementationBinding,
) (violations []checks.Violation) {
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

		implementationInterfaceMethods := make([]methodDeclaration, 0, len(interfaceMethodNames))
		for _, method := range implementationMethods {
			if interfaceMethodNamesSet[method.Name] {
				implementationInterfaceMethods = append(implementationInterfaceMethods, method)
			}
		}

		if len(implementationInterfaceMethods) != len(interfaceMethodNames) {
			violations = append(violations, checks.Violation{
				Position: binding.Position,
				Rule:     checks.DiagnosticImplementationOrder,
				Message: fmt.Sprintf(
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

			violations = append(violations, checks.Violation{
				Position: implementationInterfaceMethods[index].Position,
				Rule:     checks.DiagnosticImplementationOrder,
				Message: fmt.Sprintf(
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

/* --------------------------------------- Mock Resolution -------------------------------------- */

func resolveMockMethodsForInterface(
	interfaceName string,
	mocks map[string][]methodDeclaration,
) (
	methods []methodDeclaration,
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
