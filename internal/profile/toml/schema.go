package toml

import "ciphera/tools/internal/policy"

type schemaConfig struct {
	SchemaVersion int `toml:"schema_version"`

	Repository schemaRepositoryConfig `toml:"repository"`
	StyleGuide schemaStyleGuideConfig `toml:"style_guide"`

	PathRoles map[string][]string            `toml:"path_roles"`
	FileSets  map[string]schemaFileSetConfig `toml:"file_sets"`

	Tools   map[string]schemaPinnedTool `toml:"tools"`
	Targets map[string]schemaTarget     `toml:"targets"`

	Packs map[string]any      `toml:"packs"`
	Rules []schemaRuleBinding `toml:"rules"`
}

func decodeConfig(schema schemaConfig) (config policy.Config, err error) {
	enabledPacks, err := decodeEnabledPacks(schema.Packs)
	if err != nil {
		return policy.Config{}, err
	}

	return policy.Config{
		SchemaVersion: schema.SchemaVersion,

		Repository: decodeRepository(schema.Repository),
		StyleGuide: decodeStyleGuide(schema.StyleGuide),

		PathRoles: decodePathRoles(schema.PathRoles),
		FileSets:  decodeFileSets(schema.FileSets),

		Tools:   decodeTools(schema.Tools),
		Targets: decodeTargets(schema.Targets),

		EnabledPacks: enabledPacks,
		PackConfigs:  decodePackConfigs(schema.Packs),
		Rules:        decodeRules(schema.Rules),
	}, nil
}

func encodeConfig(config policy.Config) (schema schemaConfig) {
	return schemaConfig{
		SchemaVersion: config.SchemaVersion,

		Repository: encodeRepository(config.Repository),
		StyleGuide: encodeStyleGuide(config.StyleGuide),

		PathRoles: encodePathRoles(config.PathRoles),
		FileSets:  encodeFileSets(config.FileSets),

		Tools:   encodeTools(config.Tools),
		Targets: encodeTargets(config.Targets),

		Packs: encodePacks(config.EnabledPacks, config.PackConfigs),
		Rules: encodeRules(config.Rules),
	}
}
