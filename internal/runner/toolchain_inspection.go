package runner

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runtime"
)

func InspectToolchain(
	tools []contract.Tool,
	toolIDs []string,
	environment map[string]string,
) (statuses []runtime.ToolStatus, allValid bool) {
	statuses = runtime.InspectToolsWithEnvironment(tools, toolIDs, environment)
	statusIndex := runtime.StatusesByID(statuses)
	allValid = runtime.AllToolsValid(toolIDs, statusIndex)
	return statuses, allValid
}
