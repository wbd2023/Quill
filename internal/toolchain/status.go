package toolchain

import (
	"fmt"
	"slices"
	"strings"
)

// Status represents the outcome of inspecting one tool.
type Status struct {
	Tool    Tool
	Path    string
	Version string
	Valid   bool
	Issue   string
}

// StatusMap is a tool-status lookup keyed by tool ID.
type StatusMap map[string]Status

// NewStatusMap indexes statuses by tool ID.
func NewStatusMap(statuses []Status) (m StatusMap) {
	m = make(StatusMap, len(statuses))
	for _, status := range statuses {
		m[status.Tool.ID] = status
	}
	return m
}

// AreAllValid reports whether every tool ID has a valid status.
func (m StatusMap) AreAllValid(ids []string) (valid bool) {
	for _, id := range ids {
		if status, ok := m[id]; !ok || !status.Valid {
			return false
		}
	}
	return true
}

// ExplainIssues renders one line per invalid tool, joined with newlines in tool-ID order, or the
// empty string if all the given tools are valid.
func (m StatusMap) ExplainIssues(ids []string) (message string) {
	ordered := slices.Clone(ids)
	slices.Sort(ordered)

	var lines []string
	for _, id := range ordered {
		status, ok := m[id]
		if !ok || status.Valid {
			continue
		}

		name := status.Tool.Name
		if status.Version != "" {
			lines = append(lines, fmt.Sprintf(
				"%s: %s (found %s)", name, status.Issue, status.Version))
		} else {
			lines = append(lines, fmt.Sprintf("%s: %s", name, status.Issue))
		}
	}

	return strings.Join(lines, "\n")
}
