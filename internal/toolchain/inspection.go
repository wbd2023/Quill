package toolchain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/style"
)

/* ----------------------------------------- Inspection ----------------------------------------- */

// InspectToolsWithEnvironment reports tool status using the provided environment.
func InspectToolsWithEnvironment(
	tools []style.Tool,
	capabilities map[string]Capability,
	toolIDs []string,
	environment map[string]string,
	runner CommandRunner,
) (statuses []Status) {
	uniqueIDs := SortedUniqueToolIDs(toolIDs)
	toolByID := toolsByID(tools)
	statuses = make([]Status, 0, len(uniqueIDs))

	for _, toolID := range uniqueIDs {
		tool, found := toolByID[toolID]
		if !found {
			statuses = append(statuses, Status{
				Tool:  style.Tool{ID: toolID, Name: toolID},
				Valid: false,
				Issue: "tool is not defined in the active Packs",
			})
			continue
		}

		capability, found := capabilities[toolID]
		if !found {
			statuses = append(statuses, Status{
				Tool:  tool,
				Valid: false,
				Issue: "tool capability is not defined in the active Packs",
			})
			continue
		}

		statuses = append(statuses, inspectTool(tool, capability, environment, runner))
	}

	return statuses
}

func toolsByID(tools []style.Tool) (indexed map[string]style.Tool) {
	indexed = make(map[string]style.Tool, len(tools))
	for _, tool := range tools {
		indexed[tool.ID] = tool
	}

	return indexed
}

/* ------------------------------------------ Detection ----------------------------------------- */

func inspectTool(
	tool style.Tool,
	capability Capability,
	environment map[string]string,
	runner CommandRunner,
) (status Status) {
	status = Status{Tool: tool}

	path, err := ResolveCommandPath(capability.Command, environment)
	if err != nil {
		status.Issue = "missing from PATH"
		return status
	}

	status.Path = path
	version, versionErr := detectVersion(runner, capability, path, environment)
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

// ResolveCommandPath resolves command to an executable path. It honours the provided
// environment's PATH when set, otherwise falls back to exec.LookPath. Absolute paths and paths
// containing a separator are returned as-is.
func ResolveCommandPath(command string, environment map[string]string) (path string, err error) {
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
