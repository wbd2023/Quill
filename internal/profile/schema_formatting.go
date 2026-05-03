package profile

type schemaFormattingConfig struct {
	SectionHeaders schemaSectionHeaderConfig `toml:"section_headers"`
}

type schemaSectionHeaderConfig struct {
	RequiredMinLines  int      `toml:"required_min_lines"`
	ShortFileMaxLines int      `toml:"short_file_max_lines"`
	OveruseThreshold  int      `toml:"overuse_threshold"`
	GenericNames      []string `toml:"generic_names"`
	StructuralNames   []string `toml:"structural_names"`
}
