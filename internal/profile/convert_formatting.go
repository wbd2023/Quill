package profile

import "ciphera/tools/internal/policy"

func formattingFromSchema(schema schemaFormattingConfig) (config policy.FormattingConfig) {
	return policy.FormattingConfig{
		SectionHeaders: policy.SectionHeaderConfig{
			RequiredMinLines:  schema.SectionHeaders.RequiredMinLines,
			ShortFileMaxLines: schema.SectionHeaders.ShortFileMaxLines,
			OveruseThreshold:  schema.SectionHeaders.OveruseThreshold,
			GenericNames:      append([]string{}, schema.SectionHeaders.GenericNames...),
			StructuralNames:   append([]string{}, schema.SectionHeaders.StructuralNames...),
		},
	}
}

func formattingToSchema(config policy.FormattingConfig) (schema schemaFormattingConfig) {
	return schemaFormattingConfig{
		SectionHeaders: schemaSectionHeaderConfig{
			RequiredMinLines:  config.SectionHeaders.RequiredMinLines,
			ShortFileMaxLines: config.SectionHeaders.ShortFileMaxLines,
			OveruseThreshold:  config.SectionHeaders.OveruseThreshold,
			GenericNames:      append([]string{}, config.SectionHeaders.GenericNames...),
			StructuralNames:   append([]string{}, config.SectionHeaders.StructuralNames...),
		},
	}
}
