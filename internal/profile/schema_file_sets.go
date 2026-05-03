package profile

type schemaFileSetConfig struct {
	Name                 string              `toml:"name"`
	Extensions           []string            `toml:"extensions"`
	ExplicitFiles        map[string][]string `toml:"explicit_files"`
	PathPrefixes         map[string][]string `toml:"path_prefixes"`
	ExcludedExtensions   []string            `toml:"excluded_extensions"`
	ExcludedNames        []string            `toml:"excluded_names"`
	ExcludedNamePrefixes []string            `toml:"excluded_name_prefixes"`
	SkipBinary           bool                `toml:"skip_binary"`
}
