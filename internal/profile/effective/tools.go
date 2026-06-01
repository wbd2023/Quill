package effective

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func pinTools(
	config policy.Config,
	definitions []contract.Tool,
	availableTools map[string]contract.Tool,
) (tools []contract.Tool, err error) {
	for _, pinnedTool := range config.Tools {
		if _, found := availableTools[pinnedTool.ID]; !found {
			return nil, fmt.Errorf(
				"pinned tool %q does not match an active tool definition",
				pinnedTool.ID,
			)
		}
	}

	tools = make([]contract.Tool, 0, len(definitions))
	for _, definition := range definitions {
		pinnedTool, found := config.Tools.Lookup(definition.ID)
		if !found {
			return nil, fmt.Errorf("active tool %q is missing a pinned tool", definition.ID)
		}

		definition.PinnedVersion = pinnedTool.Version
		definition.TimeoutSeconds = pinnedTool.TimeoutSeconds
		definition.OutputLimitBytes = pinnedTool.OutputLimitBytes
		tools = append(tools, definition)
	}

	return tools, nil
}
