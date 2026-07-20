package relationships

import (
	"fmt"
	"strings"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

func ambiguousMockViolation(
	interfaceName string,
	interfaceDeclaration interfaceDeclaration,
	mockMethods []methodDeclaration,
	ambiguousMockNames []string,
) (violation analysis.Violation) {
	position := interfaceDeclaration.Position
	if len(mockMethods) > 0 {
		position = mockMethods[0].Position
	}

	return analysis.Violation{
		Position: position,
		Rule:     analysis.DiagnosticMockOrder,
		Message: fmt.Sprintf(
			"multiple mock types match interface %q: %s",
			interfaceName,
			strings.Join(ambiguousMockNames, ", "),
		),
	}
}

func mockMethodCountViolation(
	interfaceName string,
	mockName string,
	interfaceDeclaration interfaceDeclaration,
	mockMethods []methodDeclaration,
	interfaceMethodNames []string,
	mockMethodNames []string,
) (violation analysis.Violation) {
	position := interfaceDeclaration.Position
	if len(mockMethods) > 0 {
		position = mockMethods[0].Position
	}

	return analysis.Violation{
		Position: position,
		Rule:     analysis.DiagnosticMockOrder,
		Message: fmt.Sprintf(
			"mock %q for interface %q method count (%d) does not match interface (%d)",
			mockName,
			interfaceName,
			len(mockMethodNames),
			len(interfaceMethodNames),
		),
	}
}

func mockMethodOrderViolations(
	interfaceName string,
	mockName string,
	mockMethods []methodDeclaration,
	interfaceMethodNames []string,
	mockMethodNames []string,
) (violations []analysis.Violation) {
	for index := range interfaceMethodNames {
		if interfaceMethodNames[index] == mockMethodNames[index] {
			continue
		}

		return []analysis.Violation{{
			Position: mockMethods[index].Position,
			Rule:     analysis.DiagnosticMockOrder,
			Message: fmt.Sprintf(
				"mock %q for interface %q method order mismatch at position %d: "+
					"got %q, want %q",
				mockName,
				interfaceName,
				index+1,
				mockMethodNames[index],
				interfaceMethodNames[index],
			),
		}}
	}

	return nil
}
