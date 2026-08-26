package drivers

import (
	"context"
	"fmt"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

// ToolchainDriver checks that pinned tools are installed and healthy.
func ToolchainDriver(
	_ context.Context,
	_ execution.RunContext,
	_ style.Rule,
	job style.Job,
	toolStatuses toolchain.StatusMap,
) (result style.ExecutionResult, err error) {
	check, ok := job.(style.ToolchainCheck)
	if !ok {
		return style.ExecutionResult{}, fmt.Errorf(
			"toolchain driver received unsupported execution job %T",
			job,
		)
	}

	diagnostics := make([]style.Diagnostic, 0, len(check.ToolIDs))
	hasFailure := false
	for _, toolID := range check.ToolIDs {
		status, found := toolStatuses[toolID]
		if !found {
			hasFailure = true
			diagnostics = append(diagnostics, style.Diagnostic{
				Code:    "toolchain/invalid",
				Message: fmt.Sprintf("%s: no inspection status", toolID),
			})
			continue
		}
		if status.Valid {
			continue
		}

		hasFailure = true
		diagnostics = append(diagnostics, style.Diagnostic{
			Code:    "toolchain/invalid",
			Message: toolStatuses.ExplainIssues([]string{toolID}),
		})
	}

	if !hasFailure {
		return style.ExecutionResult{}, nil
	}

	return style.ExecutionResult{
		Diagnostics: diagnostics,
	}, nil
}
