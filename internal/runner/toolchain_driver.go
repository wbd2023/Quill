package runner

import (
	"fmt"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func ToolchainDriver(
	_ Context,
	spec style.ExecutionSpec,
	toolStatuses map[string]toolchain.Status,
) (result style.ExecutionResult, err error) {
	execution, found := spec.ToolchainExecution()
	if !found {
		return style.ExecutionResult{}, fmt.Errorf("toolchain driver received empty spec")
	}

	diagnostics := make([]style.Diagnostic, 0, len(execution.ToolIDs))
	foundFailure := false
	for _, toolID := range execution.ToolIDs {
		status, found := toolStatuses[toolID]
		if !found || status.Valid {
			continue
		}

		foundFailure = true
		diagnostics = append(diagnostics, style.Diagnostic{
			Code:    "toolchain/invalid",
			Message: toolchain.ExplainToolIssues([]string{toolID}, toolStatuses),
		})
	}

	if !foundFailure {
		return style.ExecutionResult{}, nil
	}

	return style.ExecutionResult{Diagnostics: diagnostics}, errRuleViolation
}
