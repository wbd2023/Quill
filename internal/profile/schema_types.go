package profile

type schemaConfig struct {
	SchemaVersion int                    `toml:"profile_version"`
	RulePacks     schemaRulePackConfig   `toml:"rule_packs"`
	Repository    schemaRepositoryConfig `toml:"repository"`
	StyleGuide    schemaStyleGuideConfig `toml:"styleguide"`
	Formatting    schemaFormattingConfig `toml:"formatting"`
	Imports       schemaImportsConfig    `toml:"imports"`
	Paths         map[string][]string    `toml:"paths"`
	FileSets      []schemaFileSetConfig  `toml:"file_sets"`
	Language      schemaLanguageConfig   `toml:"language"`
	Tools         []schemaToolPin        `toml:"tools"`
	Naming        schemaNamingConfig     `toml:"naming"`
	ControlPlane  schemaControlPlane     `toml:"control_plane"`
	Architecture  schemaArchitecture     `toml:"architecture"`
	Rules         []schemaRuleBinding    `toml:"rules"`
}

type schemaRulePackConfig struct {
	Enabled []string `toml:"enabled"`
}
