package execution

import (
	"sort"

	"ciphera/tools/internal/style"
)

// ToolIDsForRules tool i ds for rules.
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

// ToolIDsForFixes tool i ds for fixes.
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
