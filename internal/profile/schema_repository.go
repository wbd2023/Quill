package profile

type schemaRepositoryConfig struct {
	RootMarkers         []string            `toml:"root_markers"`
	DefaultScope        string              `toml:"default_scope"`
	Scopes              map[string][]string `toml:"scopes"`
	GlobalExclusions    []string            `toml:"global_exclusions"`
	GeneratedMarker     string              `toml:"generated_marker"`
	GeneratedProbeLimit int                 `toml:"generated_probe_limit"`
}

type schemaStyleGuideConfig struct {
	Path                string `toml:"path"`
	RequirementIDFormat string `toml:"requirement_id_format"`
}
