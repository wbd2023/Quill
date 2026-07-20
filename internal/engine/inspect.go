package engine

import (
	"context"

	"github.com/wbd2023/Quill/internal/toolchain"
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
	context, _, err := engine.prepareRunnerContext(operationContext, "")
	if err != nil {
		return ToolchainInspection{}, err
	}

	return engine.inspectTools(
		operationContext,
		context.Tools,
		toolIDs(context.Tools),
		context.ToolEnvironment,
	), nil
}

func (engine *Engine) inspectTools(
	ctx context.Context,
	tools map[string]toolchain.Tool,
	toolIDs []string,
	environment map[string]string,
) (inspection ToolchainInspection) {
	selected := selectTools(tools, toolIDs)
	statuses := toolchain.InspectTools(ctx, engine.commandRunner, selected, environment)
	return ToolchainInspection{
		Statuses: statuses,
		AllValid: toolchain.NewStatusMap(statuses).AreAllValid(toolIDs),
	}
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
