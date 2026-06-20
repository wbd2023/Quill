package toolchain

import (
	"fmt"
	"sort"
	"strings"

	"ciphera/tools/internal/style"
)

// Status is status.
type Status struct {
	Tool    style.Tool
	Path    string
	Version string
	Valid   bool
	Issue   string
}

// StatusesByID indexes tool statuses by tool ID.
func StatusesByID(statuses []Status) (indexed map[string]Status) {
	indexed = make(map[string]Status, len(statuses))
	for _, status := range statuses {
		indexed[status.Tool.ID] = status
	}

	return indexed
}

// AreAllToolsValid are all tools valid.
func AreAllToolsValid(toolIDs []string, indexed map[string]Status) (valid bool) {
	for _, toolID := range SortedUniqueToolIDs(toolIDs) {
		status, found := indexed[toolID]
		if !found || !status.Valid {
			return false
		}
	}

	return true
}

// ExplainToolIssues explain tool issues.
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

// SortedUniqueToolIDs sorted unique tool i ds.
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
