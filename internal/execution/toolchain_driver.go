package execution

import (
	"context"
	"fmt"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

// ToolchainDriver checks that pinned tools are installed and healthy.
func ToolchainDriver(
	_ context.Context,
	_ RunContext,
	job style.Job,
	toolStatuses toolchain.StatusMap,
) (result style.ExecutionResult, err error) {
	execution, found := job.(style.ToolchainExecution)
	if !found {
		return style.ExecutionResult{}, fmt.Errorf("toolchain driver received wrong job type")
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
