package policy

import "ciphera/tools/internal/contract"

// LanguageConfig defines the language backends available to rule bindings.
type LanguageConfig struct {
	Backends []LanguageBackendConfig
}

// LanguageBackendConfig binds a language executor to a repository scope and workdir-relative paths.
type LanguageBackendConfig struct {
	Name        string
	Language    string
	Scope       contract.Scope
	WorkDir     string
	FormatPaths []string
	CheckPaths  []string
}

// LookupBackend returns the named language backend.
func (l LanguageConfig) LookupBackend(name string) (backend LanguageBackendConfig, found bool) {
	for _, candidate := range l.Backends {
		if candidate.Name == name {
			return candidate, true
		}
	}

	return LanguageBackendConfig{}, false
}
