package policy

import "ciphera/tools/internal/requirementid"

// SchemaVersion is the current style profile schema version.
const SchemaVersion = 1

// Config is a typed style profile.
type Config struct {
	SchemaVersion  int
	Repository     RepositoryConfig
	StyleGuide     StyleGuideConfig
	Paths          PathClasses
	FileSets       FileSets
	Language       LanguageConfig
	Go             GoConfig
	Tools          PinnedTools
	Formatting     FormattingConfig
	Vocabulary     VocabularyConfig
	QualitySurface QualitySurfaceConfig
	RulePacks      RulePackConfig
	Rules          []RuleBinding
}

// StyleGuideConfig describes how the style guide is located and how its requirement IDs are parsed.
type StyleGuideConfig struct {
	Path                string
	RequirementIDScheme requirementid.Scheme
}
