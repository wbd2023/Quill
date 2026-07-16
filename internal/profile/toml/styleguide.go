package toml

import "ciphera/tools/internal/policy"

type schemaStyleGuideConfig struct {
	Path string `toml:"path"`
}

func decodeStyleGuide(schema schemaStyleGuideConfig) (styleGuide policy.StyleGuideConfig) {
	return policy.StyleGuideConfig{
		Path: schema.Path,
	}
}

func encodeStyleGuide(styleGuide policy.StyleGuideConfig) (schema schemaStyleGuideConfig) {
	return schemaStyleGuideConfig{
		Path: styleGuide.Path,
	}
}
