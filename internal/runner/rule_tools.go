package runner

import (
	"sort"

	"ciphera/tools/internal/contract"
)

func ToolIDsForRules(rules []contract.Rule) (toolIDs []string) {
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

func ToolIDsForFixes(rules []contract.Rule) (toolIDs []string) {
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
