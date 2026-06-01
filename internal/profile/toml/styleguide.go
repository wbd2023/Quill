package toml

import (
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/requirementid"
)

type schemaStyleGuideConfig struct {
	Path     string `toml:"path"`
	IDScheme string `toml:"id_scheme"`
}

func decodeStyleGuide(schema schemaStyleGuideConfig) (styleGuide policy.StyleGuideConfig) {
	return policy.StyleGuideConfig{
		Path:     schema.Path,
		IDScheme: requirementid.Scheme(schema.IDScheme),
	}
}

func encodeStyleGuide(styleGuide policy.StyleGuideConfig) (schema schemaStyleGuideConfig) {
	return schemaStyleGuideConfig{
		Path:     styleGuide.Path,
		IDScheme: string(styleGuide.IDScheme),
	}
}
