package cli

import (
	"slices"

	"ciphera/tools/internal/toolchain"
)

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
