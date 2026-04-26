package profile

type schemaFormattingConfig struct {
	SectionHeaders schemaSectionHeaderConfig `toml:"section_headers"`
}

type schemaSectionHeaderConfig struct {
	RequiredMinLines  int      `toml:"required_min_lines"`
	ShortFileMaxLines int      `toml:"short_file_max_lines"`
	OveruseCount      int      `toml:"overuse_header_count"`
	GenericNames      []string `toml:"generic_names"`
	StructuralNames   []string `toml:"structural_names"`
}

type schemaFileSetConfig struct {
	Name                 string              `toml:"name"`
	Extensions           []string            `toml:"extensions"`
	Files                map[string][]string `toml:"files"`
	Prefixes             map[string][]string `toml:"prefixes"`
	ExcludedExtensions   []string            `toml:"excluded_extensions"`
	ExcludedNames        []string            `toml:"excluded_names"`
	ExcludedNamePrefixes []string            `toml:"excluded_name_prefixes"`
	SkipBinary           bool                `toml:"skip_binary"`
}
