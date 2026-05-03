package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func languageFromSchema(schema schemaLanguageConfig) (language policy.LanguageConfig) {
	language.Backends = make([]policy.LanguageBackendConfig, 0, len(schema.Backends))
	for _, backend := range schema.Backends {
		language.Backends = append(language.Backends, policy.LanguageBackendConfig{
			Name:        backend.Name,
			Language:    backend.Language,
			Scope:       contract.Scope(backend.Scope),
			WorkDir:     backend.WorkDir,
			FormatPaths: append([]string{}, backend.FormatPaths...),
			CheckPaths:  append([]string{}, backend.CheckPaths...),
		})
	}

	return language
}

func languageToSchema(language policy.LanguageConfig) (schema schemaLanguageConfig) {
	schema.Backends = make([]schemaLanguageBackend, 0, len(language.Backends))
	for _, backend := range language.Backends {
		schema.Backends = append(schema.Backends, schemaLanguageBackend{
			Name:        backend.Name,
			Language:    backend.Language,
			Scope:       string(backend.Scope),
			WorkDir:     backend.WorkDir,
			FormatPaths: append([]string{}, backend.FormatPaths...),
			CheckPaths:  append([]string{}, backend.CheckPaths...),
		})
	}

	return schema
}
