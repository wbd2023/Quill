package runner

import (
	"sort"

	"ciphera/tools/internal/style"
)

func ToolIDsForRules(rules []style.Rule) (toolIDs []string) {
	seen := make(map[string]bool)
	for _, rule := range rules {
		for _, toolID := range rule.CheckToolIDs() {
			if seen[toolID] {
				continue
			}

			seen[toolID] = true
			toolIDs = append(toolIDs, toolID)
		}
	}

	sort.Strings(toolIDs)
	return toolIDs
}

func ToolIDsForFixes(rules []style.Rule) (toolIDs []string) {
	seen := make(map[string]bool)
	for _, rule := range rules {
		for _, toolID := range rule.FixToolIDs() {
			if seen[toolID] {
				continue
			}

			seen[toolID] = true
			toolIDs = append(toolIDs, toolID)
		}
	}

	sort.Strings(toolIDs)
	return toolIDs
}
