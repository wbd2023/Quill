package profile

import "ciphera/tools/internal/policy"

func fileSetsFromSchema(schemas []schemaFileSetConfig) (fileSets policy.FileSets) {
	fileSets = make(policy.FileSets, 0, len(schemas))
	for _, fileSet := range schemas {
		fileSets = append(fileSets, policy.FileSetConfig{
			Name:                 fileSet.Name,
			Extensions:           append([]string{}, fileSet.Extensions...),
			ExplicitFiles:        scopeMapFromSchema(fileSet.ExplicitFiles),
			PathPrefixes:         scopeMapFromSchema(fileSet.PathPrefixes),
			ExcludedExtensions:   append([]string{}, fileSet.ExcludedExtensions...),
			ExcludedNames:        append([]string{}, fileSet.ExcludedNames...),
			ExcludedNamePrefixes: append([]string{}, fileSet.ExcludedNamePrefixes...),
			SkipBinary:           fileSet.SkipBinary,
		})
	}

	return fileSets
}

func fileSetsToSchema(fileSets policy.FileSets) (schemas []schemaFileSetConfig) {
	schemas = make([]schemaFileSetConfig, 0, len(fileSets))
	for _, fileSet := range fileSets {
		schemas = append(schemas, schemaFileSetConfig{
			Name:                 fileSet.Name,
			Extensions:           append([]string{}, fileSet.Extensions...),
			ExplicitFiles:        scopeMapToSchema(fileSet.ExplicitFiles),
			PathPrefixes:         scopeMapToSchema(fileSet.PathPrefixes),
			ExcludedExtensions:   append([]string{}, fileSet.ExcludedExtensions...),
			ExcludedNames:        append([]string{}, fileSet.ExcludedNames...),
			ExcludedNamePrefixes: append([]string{}, fileSet.ExcludedNamePrefixes...),
			SkipBinary:           fileSet.SkipBinary,
		})
	}

	return schemas
}
