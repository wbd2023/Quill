package runner

import (
	"fmt"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

// ToolchainDriver toolchain driver.
func ToolchainDriver(
	_ Context,
	spec style.ExecutionSpec,
	toolStatuses toolchain.StatusMap,
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
			Message: toolStatuses.ExplainIssues([]string{toolID}),
		})
	}

	if !foundFailure {
		return style.ExecutionResult{}, nil
	}

	return style.ExecutionResult{Diagnostics: diagnostics}, nil
}
