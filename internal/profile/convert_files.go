package profile

import "ciphera/tools/internal/policy"

func formattingFromSchema(schema schemaFormattingConfig) (config policy.FormattingConfig) {
	return policy.FormattingConfig{
		SectionHeaders: policy.SectionHeaderConfig{
			RequiredMinLines:  schema.SectionHeaders.RequiredMinLines,
			ShortFileMaxLines: schema.SectionHeaders.ShortFileMaxLines,
			OveruseCount:      schema.SectionHeaders.OveruseCount,
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
			OveruseCount:      config.SectionHeaders.OveruseCount,
			GenericNames:      append([]string{}, config.SectionHeaders.GenericNames...),
			StructuralNames:   append([]string{}, config.SectionHeaders.StructuralNames...),
		},
	}
}

func fileSetsFromSchema(schemas []schemaFileSetConfig) (fileSets []policy.FileSetConfig) {
	fileSets = make([]policy.FileSetConfig, 0, len(schemas))
	for _, fileSet := range schemas {
		fileSets = append(fileSets, policy.FileSetConfig{
			Name:                 fileSet.Name,
			Extensions:           append([]string{}, fileSet.Extensions...),
			Files:                scopeMapFromSchema(fileSet.Files),
			Prefixes:             scopeMapFromSchema(fileSet.Prefixes),
			ExcludedExtensions:   append([]string{}, fileSet.ExcludedExtensions...),
			ExcludedNames:        append([]string{}, fileSet.ExcludedNames...),
			ExcludedNamePrefixes: append([]string{}, fileSet.ExcludedNamePrefixes...),
			SkipBinary:           fileSet.SkipBinary,
		})
	}

	return fileSets
}

func fileSetsToSchema(fileSets []policy.FileSetConfig) (schemas []schemaFileSetConfig) {
	schemas = make([]schemaFileSetConfig, 0, len(fileSets))
	for _, fileSet := range fileSets {
		schemas = append(schemas, schemaFileSetConfig{
			Name:                 fileSet.Name,
			Extensions:           append([]string{}, fileSet.Extensions...),
			Files:                scopeMapToSchema(fileSet.Files),
			Prefixes:             scopeMapToSchema(fileSet.Prefixes),
			ExcludedExtensions:   append([]string{}, fileSet.ExcludedExtensions...),
			ExcludedNames:        append([]string{}, fileSet.ExcludedNames...),
			ExcludedNamePrefixes: append([]string{}, fileSet.ExcludedNamePrefixes...),
			SkipBinary:           fileSet.SkipBinary,
		})
	}

	return schemas
}
