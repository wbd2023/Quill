package order

import (
	"sort"
	"strings"
)

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
