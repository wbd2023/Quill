package engine

import (
	"slices"

	"github.com/wbd2023/Quill/internal/toolchain"
)

func sortedTools(tools map[string]toolchain.Tool) (sorted []toolchain.Tool) {
	ids := toolIDs(tools)
	slices.Sort(ids)

	sorted = make([]toolchain.Tool, 0, len(ids))
	for _, id := range ids {
		sorted = append(sorted, tools[id])
	}
	return sorted
}

func toolIDs(tools map[string]toolchain.Tool) (ids []string) {
	ids = make([]string, 0, len(tools))
	for id := range tools {
		ids = append(ids, id)
	}
	return ids
}
