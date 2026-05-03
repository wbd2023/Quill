package profile

type schemaLanguageConfig struct {
	Backends []schemaLanguageBackend `toml:"backends"`
}

type schemaLanguageBackend struct {
	Name        string   `toml:"name"`
	Language    string   `toml:"language"`
	Scope       string   `toml:"scope"`
	WorkDir     string   `toml:"workdir"`
	FormatPaths []string `toml:"format_paths"`
	CheckPaths  []string `toml:"check_paths"`
}
