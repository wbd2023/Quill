package relationships

import (
	"sort"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

func checkMockOrderAgainstInterfaces(
	interfaces map[string]interfaceDeclaration,
	mocks map[string][]methodDeclaration,
) (violations []analysis.Violation) {
	interfaceNames := make([]string, 0, len(interfaces))
	for name := range interfaces {
		interfaceNames = append(interfaceNames, name)
	}
	sort.Strings(interfaceNames)

	for _, interfaceName := range interfaceNames {
		violations = append(violations, checkMockOrderForInterface(
			interfaceName,
			interfaces[interfaceName],
			mocks,
		)...)
	}

	return violations
}

func checkMockOrderForInterface(
	interfaceName string,
	interfaceDeclaration interfaceDeclaration,
	mocks map[string][]methodDeclaration,
) (violations []analysis.Violation) {
	mockMethods, matchedMockName, ambiguousMockNames, found := resolveMockMethodsForInterface(
		interfaceName,
		mocks,
	)
	if len(ambiguousMockNames) > 0 {
		return []analysis.Violation{ambiguousMockViolation(
			interfaceName,
			interfaceDeclaration,
			mockMethods,
			ambiguousMockNames,
		)}
	}

	if !found {
		return nil
	}

	interfaceMethodNames := methodNames(interfaceDeclaration.Methods)
	mockMethodNames := methodNames(mockMethods)
	if len(interfaceMethodNames) != len(mockMethodNames) {
		return []analysis.Violation{mockMethodCountViolation(
			interfaceName,
			matchedMockName,
			interfaceDeclaration,
			mockMethods,
			interfaceMethodNames,
			mockMethodNames,
		)}
	}

	return mockMethodOrderViolations(
		interfaceName,
		matchedMockName,
		mockMethods,
		interfaceMethodNames,
		mockMethodNames,
	)
}
