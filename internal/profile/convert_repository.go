package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func repositoryFromSchema(schema schemaRepositoryConfig) (repository policy.RepositoryConfig) {
	return policy.RepositoryConfig{
		RootMarkers:         append([]string{}, schema.RootMarkers...),
		DefaultScope:        contract.Scope(schema.DefaultScope),
		Scopes:              scopeMapFromSchema(schema.Scopes),
		GlobalExclusions:    append([]string{}, schema.GlobalExclusions...),
		GeneratedMarker:     schema.GeneratedMarker,
		GeneratedProbeLimit: schema.GeneratedProbeLimit,
	}
}

func repositoryToSchema(repository policy.RepositoryConfig) (schema schemaRepositoryConfig) {
	return schemaRepositoryConfig{
		RootMarkers:         append([]string{}, repository.RootMarkers...),
		DefaultScope:        string(repository.DefaultScope),
		Scopes:              scopeMapToSchema(repository.Scopes),
		GlobalExclusions:    append([]string{}, repository.GlobalExclusions...),
		GeneratedMarker:     repository.GeneratedMarker,
		GeneratedProbeLimit: repository.GeneratedProbeLimit,
	}
}

func styleGuideFromSchema(schema schemaStyleGuideConfig) (config policy.StyleGuideConfig) {
	return policy.StyleGuideConfig{
		Path:                schema.Path,
		RequirementIDFormat: schema.RequirementIDFormat,
	}
}
