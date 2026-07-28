package policy

// SchemaVersion is the current style profile schema version.
const SchemaVersion = 1

// Profile is a typed style profile.
type Profile struct {
	SchemaVersion int

	Repository RepositoryConfig
	StyleGuide StyleGuideConfig

	PathRoles PathRoles
	FileSets  FileSets

	Tools   PinnedTools
	Targets TargetConfigs

	EnabledPacks []string
	PackConfigs  PackConfigs
	PackSources  []PackSource
	Rules        []RuleBinding
}

// PackSource is one declared local external Pack directory.
type PackSource struct {
	Path string
}

// StyleGuideConfig describes how the style guide is located.
type StyleGuideConfig struct {
	Path string
}
