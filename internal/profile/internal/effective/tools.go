package effective

import (
	"fmt"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

func pinTools(
	config policy.Config,
	definitions []style.Tool,
	availableTools map[string]style.Tool,
) (tools []style.Tool, err error) {
	for _, pinnedTool := range config.Tools {
		if _, found := availableTools[pinnedTool.ID]; !found {
			return nil, fmt.Errorf(
				"pinned tool %q does not match an active tool definition",
				pinnedTool.ID,
			)
		}
	}

	tools = make([]style.Tool, 0, len(definitions))
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
