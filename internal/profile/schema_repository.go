package profile

type schemaRepositoryConfig struct {
	RootMarkers         []string            `toml:"root_markers"`
	DefaultScope        string              `toml:"default_scope"`
	ScopeRoots          map[string][]string `toml:"scope_roots"`
	GlobalExclusions    []string            `toml:"global_exclusions"`
	GeneratedMarker     string              `toml:"generated_marker"`
	GeneratedProbeBytes int                 `toml:"generated_probe_bytes"`
}

type schemaStyleGuideConfig struct {
	Path                string `toml:"path"`
	RequirementIDScheme string `toml:"requirement_id_scheme"`
}
