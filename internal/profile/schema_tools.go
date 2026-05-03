package profile

type schemaPinnedTool struct {
	ID               string `toml:"id"`
	Version          string `toml:"version"`
	TimeoutSeconds   int    `toml:"timeout_seconds"`
	OutputLimitBytes int64  `toml:"output_limit_bytes"`
}
