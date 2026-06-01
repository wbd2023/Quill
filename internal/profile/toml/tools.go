package toml

import "ciphera/tools/internal/policy"

type schemaPinnedTool struct {
	Version          string `toml:"version"`
	TimeoutSeconds   int    `toml:"timeout_seconds"`
	OutputLimitBytes int64  `toml:"output_limit_bytes"`
}

func decodeTools(schemas map[string]schemaPinnedTool) (tools policy.PinnedTools) {
	tools = make(policy.PinnedTools, 0, len(schemas))
	for _, id := range sortedMapKeys(schemas) {
		tool := schemas[id]
		tools = append(tools, policy.PinnedTool{
			ID:               id,
			Version:          tool.Version,
			TimeoutSeconds:   tool.TimeoutSeconds,
			OutputLimitBytes: tool.OutputLimitBytes,
		})
	}

	return tools
}

func encodeTools(tools policy.PinnedTools) (schemas map[string]schemaPinnedTool) {
	if tools == nil {
		return nil
	}

	schemas = make(map[string]schemaPinnedTool, len(tools))
	for _, tool := range tools {
		schemas[tool.ID] = schemaPinnedTool{
			Version:          tool.Version,
			TimeoutSeconds:   tool.TimeoutSeconds,
			OutputLimitBytes: tool.OutputLimitBytes,
		}
	}

	return schemas
}
