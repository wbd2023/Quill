package toolchain

import (
	"fmt"
	"sort"
	"strings"

	"ciphera/tools/internal/contract"
)

type Status struct {
	Tool    contract.Tool
	Path    string
	Version string
	Valid   bool
	Issue   string
}

func StatusesByID(statuses []Status) (indexed map[string]Status) {
	indexed = make(map[string]Status, len(statuses))
	for _, status := range statuses {
		indexed[status.Tool.ID] = status
	}

	return indexed
}

func AllToolsValid(toolIDs []string, indexed map[string]Status) (valid bool) {
	for _, toolID := range SortedUniqueToolIDs(toolIDs) {
		status, found := indexed[toolID]
		if !found || !status.Valid {
			return false
		}
	}

	return true
}

func ExplainToolIssues(toolIDs []string, indexed map[string]Status) (message string) {
	var parts []string
	for _, toolID := range SortedUniqueToolIDs(toolIDs) {
		status, found := indexed[toolID]
		if !found || status.Valid {
			continue
		}

		parts = append(parts, formatToolStatusLine(status))
	}

	return strings.Join(parts, "\n")
}

func SortedUniqueToolIDs(toolIDs []string) (deduped []string) {
	seen := make(map[string]bool)
	for _, toolID := range toolIDs {
		if seen[toolID] {
			continue
		}

		seen[toolID] = true
		deduped = append(deduped, toolID)
	}

	sort.Strings(deduped)
	return deduped
}

func formatToolStatusLine(status Status) (line string) {
	name := status.Tool.Name
	if status.Version != "" {
		return fmt.Sprintf("%s: %s (found %s)", name, status.Issue, status.Version)
	}

	return fmt.Sprintf("%s: %s", name, status.Issue)
}
