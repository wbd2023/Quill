package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"ciphera/tools/internal/contract"
)

/* -------------------------------------------- Types ------------------------------------------- */

// ToolStatus describes the detected state of one configured tool.
type ToolStatus struct {
	Tool    contract.Tool
	Path    string
	Version string
	Valid   bool
	Issue   string
}

/* ----------------------------------------- Inspection ----------------------------------------- */

// InspectTools reports tool status using the current process environment.
func InspectTools(tools []contract.Tool, toolIDs []string) (statuses []ToolStatus) {
	return InspectToolsWithEnvironment(tools, toolIDs, nil)
}

// InspectToolsWithEnvironment reports tool status using the provided environment.
func InspectToolsWithEnvironment(
	tools []contract.Tool,
	toolIDs []string,
	environment map[string]string,
) (statuses []ToolStatus) {
	uniqueIDs := sortedUniqueToolIDs(toolIDs)
	toolByID := toolsByID(tools)
	statuses = make([]ToolStatus, 0, len(uniqueIDs))
	for _, toolID := range uniqueIDs {
		tool, found := toolByID[toolID]
		if !found {
			statuses = append(statuses, ToolStatus{
				Tool:  contract.Tool{ID: toolID, Name: toolID},
				Valid: false,
				Issue: "tool is not defined in the active rule packs",
			})
			continue
		}

		statuses = append(statuses, inspectTool(tool, environment))
	}

	return statuses
}

func toolsByID(tools []contract.Tool) (indexed map[string]contract.Tool) {
	indexed = make(map[string]contract.Tool, len(tools))
	for _, tool := range tools {
		indexed[tool.ID] = tool
	}

	return indexed
}

// StatusesByID indexes tool statuses by tool identifier.
func StatusesByID(statuses []ToolStatus) (indexed map[string]ToolStatus) {
	indexed = make(map[string]ToolStatus, len(statuses))
	for _, status := range statuses {
		indexed[status.Tool.ID] = status
	}

	return indexed
}

// AllToolsValid reports whether every requested tool is present and pinned correctly.
func AllToolsValid(toolIDs []string, indexed map[string]ToolStatus) (valid bool) {
	for _, toolID := range sortedUniqueToolIDs(toolIDs) {
		status, found := indexed[toolID]
		if !found || !status.Valid {
			return false
		}
	}

	return true
}

// ExplainToolIssues renders invalid tool statuses as human-readable lines.
func ExplainToolIssues(toolIDs []string, indexed map[string]ToolStatus) (message string) {
	var parts []string
	for _, toolID := range sortedUniqueToolIDs(toolIDs) {
		status, found := indexed[toolID]
		if !found || status.Valid {
			continue
		}

		parts = append(parts, formatStatusLine(status))
	}

	return strings.Join(parts, "\n")
}

func formatStatusLine(status ToolStatus) (line string) {
	name := status.Tool.Name
	if status.Version != "" {
		return fmt.Sprintf("%s: %s (found %s)", name, status.Issue, status.Version)
	}

	return fmt.Sprintf("%s: %s", name, status.Issue)
}

/* ------------------------------------------ Detection ----------------------------------------- */

func inspectTool(tool contract.Tool, environment map[string]string) (status ToolStatus) {
	status = ToolStatus{Tool: tool}

	path, err := lookCommandPath(tool.Command, environment)
	if err != nil {
		status.Issue = "missing from PATH"
		return status
	}

	status.Path = path
	version, versionErr := detectVersion(tool, path, environment)
	if versionErr != nil {
		status.Issue = versionErr.Error()
		return status
	}

	status.Version = version
	if !matchesPinnedVersion(version, tool.PinnedVersion) {
		status.Issue = fmt.Sprintf("requires pinned version %s", tool.PinnedVersion)
		return status
	}

	status.Valid = true
	return status
}

/* --------------------------------------- Command Lookup --------------------------------------- */

func lookCommandPath(command string, environment map[string]string) (path string, err error) {
	commandPath := lookupEnvironmentValue(environment, "PATH")
	if commandPath == "" {
		return exec.LookPath(command)
	}

	if filepath.IsAbs(command) || strings.ContainsRune(command, os.PathSeparator) {
		return command, nil
	}

	for _, directory := range filepath.SplitList(commandPath) {
		candidate := filepath.Join(directory, command)
		info, statErr := os.Stat(candidate)
		if statErr != nil || info.IsDir() || !isExecutable(info.Mode()) {
			continue
		}

		return candidate, nil
	}

	return "", exec.ErrNotFound
}

/* ------------------------------------- Environment Lookup ------------------------------------- */

func lookupEnvironmentValue(environment map[string]string, key string) (value string) {
	if environment != nil {
		if value, found := environment[key]; found {
			return value
		}
	}

	return os.Getenv(key)
}

/* -------------------------------------- Executable Checks ------------------------------------- */

func isExecutable(mode os.FileMode) (executable bool) {
	return mode&0o111 != 0
}

/* ------------------------------------------ Tool IDs ------------------------------------------ */

func sortedUniqueToolIDs(toolIDs []string) (deduped []string) {
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
