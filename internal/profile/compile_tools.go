package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func bindPinnedTools(
	config policy.Config,
	tools []contract.Tool,
) (pinned []contract.Tool, err error) {
	toolByID := make(map[string]contract.Tool, len(tools))
	for _, tool := range tools {
		toolByID[tool.ID] = tool
	}

	for _, pinnedTool := range config.Tools {
		if _, found := toolByID[pinnedTool.ID]; !found {
			return nil, fmt.Errorf(
				"pinned tool %q does not match an active tool definition",
				pinnedTool.ID,
			)
		}
	}

	pinned = make([]contract.Tool, 0, len(tools))
	for _, tool := range tools {
		pinnedTool, found := config.Tools.Lookup(tool.ID)
		if !found {
			return nil, fmt.Errorf("active tool %q is missing a pinned tool", tool.ID)
		}

		tool.PinnedVersion = pinnedTool.Version
		tool.TimeoutSeconds = pinnedTool.TimeoutSeconds
		tool.OutputLimitBytes = pinnedTool.OutputLimitBytes
		pinned = append(pinned, tool)
	}

	return pinned, nil
}
