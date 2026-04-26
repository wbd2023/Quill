package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func bindToolPins(
	config policy.Config,
	tools []contract.Tool,
) (pinned []contract.Tool, err error) {
	toolByID := make(map[string]contract.Tool, len(tools))
	for _, tool := range tools {
		toolByID[tool.ID] = tool
	}

	for _, pin := range config.Tools {
		if _, found := toolByID[pin.ID]; !found {
			return nil, fmt.Errorf("tool pin %q does not match an active builtin tool", pin.ID)
		}
	}

	pinned = make([]contract.Tool, 0, len(tools))
	for _, tool := range tools {
		pin, found := config.ToolPin(tool.ID)
		if !found {
			return nil, fmt.Errorf("active tool %q is missing a tool pin", tool.ID)
		}

		tool.PinnedVersion = pin.Version
		tool.TimeoutSeconds = pin.TimeoutSeconds
		tool.OutputLimitBytes = pin.OutputLimitBytes
		pinned = append(pinned, tool)
	}

	return pinned, nil
}
