package runner

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

func ToolchainDriver(
	_ Context,
	spec contract.ExecutionSpec,
	toolStatuses map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.ToolchainExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf("toolchain driver received empty spec")
	}

	diagnostics := make([]contract.Diagnostic, 0, len(execution.ToolIDs))
	foundFailure := false
	for _, toolID := range execution.ToolIDs {
		status, found := toolStatuses[toolID]
		if !found || status.Valid {
			continue
		}

		foundFailure = true
		diagnostics = append(diagnostics, contract.Diagnostic{
			Code:    "toolchain/invalid",
			Message: toolchain.ExplainToolIssues([]string{toolID}, toolStatuses),
		})
	}

	if !foundFailure {
		return contract.ExecutionResult{}, nil
	}

	return contract.ExecutionResult{Diagnostics: diagnostics}, errRuleViolation
}
