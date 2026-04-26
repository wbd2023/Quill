package policy

import "ciphera/tools/internal/contract"

type Config struct {
	SchemaVersion int
	RulePacks     RulePackConfig
	Repository    RepositoryConfig
	StyleGuide    StyleGuideConfig
	Formatting    FormattingConfig
	Imports       ImportsConfig
	Paths         PathClassSet
	FileSets      []FileSetConfig
	Language      LanguageConfig
	Tools         []ToolPin
	Naming        NamingConfig
	ControlPlane  ControlPlaneConfig
	Architecture  ArchitectureConfig
	Rules         []RuleBinding
}

type RulePackConfig struct {
	Enabled []string
}

type StyleGuideConfig struct {
	Path                string
	RequirementIDFormat string
}

type ImportsConfig struct {
	LocalPrefix string
}

func (config Config) FileSet(name string) (fileSet FileSetConfig, found bool) {
	for _, fileSet := range config.FileSets {
		if fileSet.Name == name {
			return fileSet, true
		}
	}

	return FileSetConfig{}, false
}

func (config Config) LanguageBackend(
	name string,
) (backend LanguageBackendConfig, found bool) {
	for _, backend := range config.Language.Backends {
		if backend.Name == name {
			return backend, true
		}
	}

	return LanguageBackendConfig{}, false
}

func (config Config) ToolPin(id string) (pin ToolPin, found bool) {
	for _, tool := range config.Tools {
		if tool.ID == id {
			return tool, true
		}
	}

	return ToolPin{}, false
}

func (config Config) ScopeExists(scope contract.Scope) (found bool) {
	return config.Repository.ScopeExists(scope)
}
