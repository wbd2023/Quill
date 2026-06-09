package cli

import (
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func toolIDsFromTools(tools []style.Tool) (toolIDs []string) {
	for _, tool := range tools {
		toolIDs = append(toolIDs, tool.ID)
	}

	return toolIDs
}

func inspectToolchain(
	tools []style.Tool,
	capabilities map[string]toolchain.Capability,
	toolIDs []string,
	environment map[string]string,
) (statuses []toolchain.Status, allValid bool) {
	statuses = toolchain.InspectToolsWithEnvironment(tools, capabilities, toolIDs, environment, runtime.RunToolchainCommand)
	statusIndex := toolchain.StatusesByID(statuses)
	allValid = toolchain.AllToolsValid(toolIDs, statusIndex)
	return statuses, allValid
}
