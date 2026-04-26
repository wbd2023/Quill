package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

/* ----------------------------------------- Inspection ----------------------------------------- */

// InspectTools reports tool status using the current process environment.
func InspectTools(
	tools []contract.Tool,
	capabilities map[string]toolchain.Capability,
	toolIDs []string,
) (statuses []toolchain.Status) {
	return InspectToolsWithEnvironment(tools, capabilities, toolIDs, nil)
}

// InspectToolsWithEnvironment reports tool status using the provided environment.
func InspectToolsWithEnvironment(
	tools []contract.Tool,
	capabilities map[string]toolchain.Capability,
	toolIDs []string,
	environment map[string]string,
) (statuses []toolchain.Status) {
	uniqueIDs := toolchain.SortedUniqueToolIDs(toolIDs)
	toolByID := toolsByID(tools)
	statuses = make([]toolchain.Status, 0, len(uniqueIDs))

	for _, toolID := range uniqueIDs {
		tool, found := toolByID[toolID]
		if !found {
			statuses = append(statuses, toolchain.Status{
				Tool:  contract.Tool{ID: toolID, Name: toolID},
				Valid: false,
				Issue: "tool is not defined in the active rule packs",
			})
			continue
		}

		capability, found := capabilities[toolID]
		if !found {
			statuses = append(statuses, toolchain.Status{
				Tool:  tool,
				Valid: false,
				Issue: "tool capability is not defined in the active rule packs",
			})
			continue
		}

		statuses = append(statuses, inspectTool(tool, capability, environment))
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

/* ------------------------------------------ Detection ----------------------------------------- */

func inspectTool(
	tool contract.Tool,
	capability toolchain.Capability,
	environment map[string]string,
) (status toolchain.Status) {
	status = toolchain.Status{Tool: tool}

	path, err := lookupCommandPath(capability.Command, environment)
	if err != nil {
		status.Issue = "missing from PATH"
		return status
	}

	status.Path = path
	version, versionErr := detectVersion(tool, capability, path, environment)
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

func lookupCommandPath(command string, environment map[string]string) (path string, err error) {
	pathList := lookupEnvironmentVariable(environment, "PATH")
	if pathList == "" {
		return exec.LookPath(command)
	}

	if filepath.IsAbs(command) || strings.ContainsRune(command, os.PathSeparator) {
		return command, nil
	}

	for _, directory := range filepath.SplitList(pathList) {
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

func lookupEnvironmentVariable(environment map[string]string, key string) (value string) {
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
