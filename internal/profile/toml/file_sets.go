package toml

import "ciphera/tools/internal/policy"

type schemaFileSetConfig struct {
	Include schemaFileSetInclude `toml:"include"`
	Exclude schemaFileSetExclude `toml:"exclude"`
}

type schemaFileSetInclude struct {
	Extensions []string            `toml:"extensions"`
	Files      map[string][]string `toml:"files"`
	Paths      map[string][]string `toml:"paths"`
}

type schemaFileSetExclude struct {
	Extensions []string `toml:"extensions"`
	Files      []string `toml:"files"`
}

func decodeFileSets(schemas map[string]schemaFileSetConfig) (fileSets policy.FileSets) {
	fileSets = make(policy.FileSets, 0, len(schemas))
	for _, name := range sortedMapKeys(schemas) {
		fileSet := schemas[name]
		fileSets = append(fileSets, policy.FileSetConfig{
			Name: name,
			Include: policy.FileSetInclude{
				Extensions: append([]string{}, fileSet.Include.Extensions...),
				Files:      decodeScopeMap(fileSet.Include.Files),
				Paths:      decodeScopeMap(fileSet.Include.Paths),
			},
			Exclude: policy.FileSetExclude{
				Extensions: append([]string{}, fileSet.Exclude.Extensions...),
				Files:      append([]string{}, fileSet.Exclude.Files...),
			},
		})
	}

	return fileSets
}

func encodeFileSets(fileSets policy.FileSets) (schemas map[string]schemaFileSetConfig) {
	if fileSets == nil {
		return nil
	}

	schemas = make(map[string]schemaFileSetConfig, len(fileSets))
	for _, fileSet := range fileSets {
		schemas[fileSet.Name] = schemaFileSetConfig{
			Include: schemaFileSetInclude{
				Extensions: append([]string{}, fileSet.Include.Extensions...),
				Files:      encodeScopeMap(fileSet.Include.Files),
				Paths:      encodeScopeMap(fileSet.Include.Paths),
			},
			Exclude: schemaFileSetExclude{
				Extensions: append([]string{}, fileSet.Exclude.Extensions...),
				Files:      append([]string{}, fileSet.Exclude.Files...),
			},
		}
	}

	return schemas
}
