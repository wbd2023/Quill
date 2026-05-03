package profile

import "ciphera/tools/internal/policy"

func toolsFromSchema(schemas []schemaPinnedTool) (tools policy.PinnedTools) {
	tools = make(policy.PinnedTools, 0, len(schemas))
	for _, tool := range schemas {
		tools = append(tools, policy.PinnedTool{
			ID:               tool.ID,
			Version:          tool.Version,
			TimeoutSeconds:   tool.TimeoutSeconds,
			OutputLimitBytes: tool.OutputLimitBytes,
		})
	}

	return tools
}

func toolsToSchema(tools policy.PinnedTools) (schemas []schemaPinnedTool) {
	schemas = make([]schemaPinnedTool, 0, len(tools))
	for _, tool := range tools {
		schemas = append(schemas, schemaPinnedTool{
			ID:               tool.ID,
			Version:          tool.Version,
			TimeoutSeconds:   tool.TimeoutSeconds,
			OutputLimitBytes: tool.OutputLimitBytes,
		})
	}

	return schemas
}
