package profile

type schemaConfig struct {
	SchemaVersion  int                    `toml:"profile_version"`
	Repository     schemaRepositoryConfig `toml:"repository"`
	StyleGuide     schemaStyleGuideConfig `toml:"styleguide"`
	Paths          map[string][]string    `toml:"paths"`
	FileSets       []schemaFileSetConfig  `toml:"file_sets"`
	Language       schemaLanguageConfig   `toml:"language"`
	Go             schemaGoConfig         `toml:"go"`
	Tools          []schemaPinnedTool     `toml:"tools"`
	Formatting     schemaFormattingConfig `toml:"formatting"`
	Vocabulary     schemaVocabularyConfig `toml:"vocabulary"`
	QualitySurface schemaQualitySurface   `toml:"quality_surface"`
	RulePacks      schemaRulePackConfig   `toml:"rule_packs"`
	Rules          []schemaRuleBinding    `toml:"rules"`
}

type schemaRulePackConfig struct {
	Enabled []string `toml:"enabled"`
}
