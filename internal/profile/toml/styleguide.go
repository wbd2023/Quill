package toml

import "github.com/wbd2023/Quill/internal/policy"

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
