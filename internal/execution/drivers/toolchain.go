package drivers

import (
	"context"
	"fmt"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

// ToolchainDriver checks that pinned tools are installed and healthy.
func ToolchainDriver(
	_ context.Context,
	_ execution.RunContext,
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
