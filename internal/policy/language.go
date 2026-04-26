package policy

import "ciphera/tools/internal/contract"

type LanguageConfig struct {
	Backends []LanguageBackendConfig
}

type LanguageBackendConfig struct {
	Name        string
	Language    string
	Scope       contract.Scope
	Workdir     string
	FormatPaths []string
	StylePaths  []string
}
