package toolchain

import (
	"fmt"
	"slices"

	"ciphera/tools/internal/runtime"
)

// InspectTools reports the status of each tool in tools, sorted by tool ID.
func InspectTools(tools map[string]Tool, environment map[string]string) (statuses []Status) {
	ids := make([]string, 0, len(tools))
	for id := range tools {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	statuses = make([]Status, 0, len(ids))

	for _, id := range ids {
		statuses = append(statuses, inspectTool(tools[id], environment))
	}

	return statuses
}

func inspectTool(tool Tool, environment map[string]string) (status Status) {
	status = Status{Tool: tool}

	path, err := runtime.ResolveCommandPath(tool.Command, environment)
	if err != nil {
		status.Issue = "missing from PATH"
		return status
	}

	status.Path = path
	version, err := tool.Version(environment, path)
	if err != nil {
		status.Issue = err.Error()
		return status
	}

	status.Version = version
	if normaliseVersion(version) != normaliseVersion(tool.PinnedVersion) {
		status.Issue = fmt.Sprintf("requires pinned version %s", tool.PinnedVersion)
		return status
	}

	status.Valid = true
	return status
}
