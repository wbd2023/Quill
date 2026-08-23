package profile

type schemaStyleGuideConfig struct {
	Path string `toml:"path"`
}

func decodeStyleGuide(schema schemaStyleGuideConfig) (styleGuide StyleGuideConfig) {
	return StyleGuideConfig(schema)
}

func encodeStyleGuide(styleGuide StyleGuideConfig) (schema schemaStyleGuideConfig) {
	return schemaStyleGuideConfig(styleGuide)
}
