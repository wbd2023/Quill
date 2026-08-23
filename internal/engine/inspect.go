package engine

import (
	"context"

	"github.com/wbd2023/quill/internal/toolchain"
)

// ToolchainInspection contains structured tool inspection outcomes.
type ToolchainInspection struct {
	Statuses []toolchain.Status
	AllValid bool
}

// Inspect loads the repository and inspects every configured tool.
func (engine *Engine) Inspect(
	operationContext context.Context,
) (inspection ToolchainInspection, operationError error) {
	runContext, _, err := engine.prepareRun(operationContext, "")
	if err != nil {
		return ToolchainInspection{}, err
	}

	inspection, err = engine.inspectTools(
		operationContext,
		runContext.Tools,
		toolIDs(runContext.Tools),
		runContext.ToolEnvironment,
	)
	if err != nil {
		return ToolchainInspection{}, err
	}

	return inspection, nil
}

func (engine *Engine) inspectTools(
	ctx context.Context,
	tools map[string]toolchain.Tool,
	toolIDs []string,
	environment map[string]string,
) (inspection ToolchainInspection, inspectionError error) {
	selected := selectTools(tools, toolIDs)
	statuses, err := toolchain.InspectTools(ctx, engine.commandRunner, selected, environment)
	if err != nil {
		return ToolchainInspection{}, err
	}

	return ToolchainInspection{
		Statuses: statuses,
		AllValid: toolchain.NewStatusMap(statuses).AreAllValid(toolIDs),
	}, nil
}

func selectTools(
	tools map[string]toolchain.Tool,
	toolIDs []string,
) (selected map[string]toolchain.Tool) {
	selected = make(map[string]toolchain.Tool, len(toolIDs))
	for _, toolID := range toolIDs {
		selected[toolID] = tools[toolID]
	}
	return selected
}
