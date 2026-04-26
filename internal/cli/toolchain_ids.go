package cli

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

func toolIDsFromTools(tools []contract.Tool) (toolIDs []string) {
	for _, tool := range tools {
		toolIDs = append(toolIDs, tool.ID)
	}

	return toolIDs
}

func inspectToolchain(
	tools []contract.Tool,
	capabilities map[string]toolchain.Capability,
	toolIDs []string,
	environment map[string]string,
) (statuses []toolchain.Status, allValid bool) {
	statuses = runtime.InspectToolsWithEnvironment(tools, capabilities, toolIDs, environment)
	statusIndex := toolchain.StatusesByID(statuses)
	allValid = toolchain.AllToolsValid(toolIDs, statusIndex)
	return statuses, allValid
}
