package relationships

import (
	"fmt"
	"sort"

	"ciphera/tools/internal/checks/golang/analysis"
)

/* ----------------------------------- Implementation Matching ---------------------------------- */

func checkImplementationOrderAgainstInterfaces(
	interfaces map[string]interfaceDeclaration,
	implementations map[string][]methodDeclaration,
	bindings []implementationBinding,
) (violations []analysis.Violation) {
	sort.Slice(bindings, func(i int, j int) bool {
		if bindings[i].InterfaceName == bindings[j].InterfaceName {
			return bindings[i].ImplementationName < bindings[j].ImplementationName
		}
		return bindings[i].InterfaceName < bindings[j].InterfaceName
	})

	for _, binding := range bindings {
		violations = append(violations, checkImplementationOrderForBinding(
			binding,
			interfaces,
			implementations,
		)...)
	}

	return violations
}

func checkImplementationOrderForBinding(
	binding implementationBinding,
	interfaces map[string]interfaceDeclaration,
	implementations map[string][]methodDeclaration,
) (violations []analysis.Violation) {
	interfaceDeclaration, found := interfaces[binding.InterfaceName]
	if !found {
		return nil
	}

	implementationMethods, found := implementations[binding.ImplementationKey]
	if !found {
		return nil
	}

	interfaceMethodNames := methodNames(interfaceDeclaration.Methods)
	implementationInterfaceMethods := matchingMethods(
		implementationMethods,
		methodNameSet(interfaceMethodNames),
	)
	if len(implementationInterfaceMethods) != len(interfaceMethodNames) {
		return []analysis.Violation{implementationMethodCountViolation(
			binding,
			implementationInterfaceMethods,
			interfaceMethodNames,
		)}
	}

	return implementationMethodOrderViolations(
		binding,
		implementationInterfaceMethods,
		interfaceMethodNames,
	)
}

/* ----------------------------------------- Diagnostics ---------------------------------------- */

func implementationMethodCountViolation(
	binding implementationBinding,
	implementationMethods []methodDeclaration,
	interfaceMethodNames []string,
) (violation analysis.Violation) {
	return analysis.Violation{
		Position: binding.Position,
		Rule:     analysis.DiagnosticImplementationOrder,
		Message: fmt.Sprintf(
			"implementation %q for interface %q method count (%d) "+
				"does not match interface (%d)",
			binding.ImplementationName,
			binding.InterfaceName,
			len(implementationMethods),
			len(interfaceMethodNames),
		),
	}
}

func implementationMethodOrderViolations(
	binding implementationBinding,
	implementationMethods []methodDeclaration,
	interfaceMethodNames []string,
) (violations []analysis.Violation) {
	for index := range interfaceMethodNames {
		if implementationMethods[index].Name == interfaceMethodNames[index] {
			continue
		}

		return []analysis.Violation{{
			Position: implementationMethods[index].Position,
			Rule:     analysis.DiagnosticImplementationOrder,
			Message: fmt.Sprintf(
				"implementation %q for interface %q method order mismatch at position %d: "+
					"got %q, want %q",
				binding.ImplementationName,
				binding.InterfaceName,
				index+1,
				implementationMethods[index].Name,
				interfaceMethodNames[index],
			),
		}}
	}

	return nil
}
