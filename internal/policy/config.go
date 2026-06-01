package policy

import "ciphera/tools/internal/requirementid"

// SchemaVersion is the current style profile schema version.
const SchemaVersion = 1

// Config is a typed style profile.
type Config struct {
	SchemaVersion int

	Repository RepositoryConfig
	StyleGuide StyleGuideConfig

	PathRoles PathRoles
	FileSets  FileSets

	Tools   PinnedTools
	Targets TargetConfigs

	EnabledPacks []string
	PackConfigs  PackConfigs
	Rules        []RuleBinding
}

// StyleGuideConfig describes how the style guide is located and how its requirement IDs are parsed.
type StyleGuideConfig struct {
	Path     string
	IDScheme requirementid.Scheme
}
