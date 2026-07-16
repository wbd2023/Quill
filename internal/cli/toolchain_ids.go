package cli

import (
	"slices"

	"ciphera/tools/internal/toolchain"
)

func inspectToolchain(
	runner toolchain.CommandRunner,
	tools map[string]toolchain.Tool,
	environment map[string]string,
) (statuses []toolchain.Status, allValid bool) {
	statuses = toolchain.InspectTools(runner, tools, environment)
	statusMap := toolchain.NewStatusMap(statuses)
	allValid = statusMap.AreAllValid(sortedToolIDs(tools))
	return statuses, allValid
}

// selectTools returns a map containing only the tools whose IDs are in wantedIDs.
func selectTools(
	tools map[string]toolchain.Tool,
	wantedIDs []string,
) (selected map[string]toolchain.Tool) {
	selected = make(map[string]toolchain.Tool, len(wantedIDs))
	for _, toolID := range wantedIDs {
		selected[toolID] = tools[toolID]
	}
	return selected
}

func sortedToolIDs(tools map[string]toolchain.Tool) (toolIDs []string) {
	toolIDs = make([]string, 0, len(tools))
	for toolID := range tools {
		toolIDs = append(toolIDs, toolID)
	}
	slices.Sort(toolIDs)
	return toolIDs
}

func sortedTools(tools map[string]toolchain.Tool) (sorted []toolchain.Tool) {
	for _, toolID := range sortedToolIDs(tools) {
		sorted = append(sorted, tools[toolID])
	}
	return sorted
}
