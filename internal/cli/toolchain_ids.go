package cli

import "ciphera/tools/internal/contract"

func toolIDsFromTools(tools []contract.Tool) (toolIDs []string) {
	for _, tool := range tools {
		toolIDs = append(toolIDs, tool.ID)
	}

	return toolIDs
}
