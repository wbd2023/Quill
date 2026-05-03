package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/requirementid"
)

func repositoryFromSchema(schema schemaRepositoryConfig) (repository policy.RepositoryConfig) {
	return policy.RepositoryConfig{
		RootMarkers:         append([]string{}, schema.RootMarkers...),
		DefaultScope:        contract.Scope(schema.DefaultScope),
		ScopeRoots:          scopeMapFromSchema(schema.ScopeRoots),
		GlobalExclusions:    append([]string{}, schema.GlobalExclusions...),
		GeneratedMarker:     schema.GeneratedMarker,
		GeneratedProbeBytes: schema.GeneratedProbeBytes,
	}
}

func repositoryToSchema(repository policy.RepositoryConfig) (schema schemaRepositoryConfig) {
	return schemaRepositoryConfig{
		RootMarkers:         append([]string{}, repository.RootMarkers...),
		DefaultScope:        string(repository.DefaultScope),
		ScopeRoots:          scopeMapToSchema(repository.ScopeRoots),
		GlobalExclusions:    append([]string{}, repository.GlobalExclusions...),
		GeneratedMarker:     repository.GeneratedMarker,
		GeneratedProbeBytes: repository.GeneratedProbeBytes,
	}
}

func styleGuideFromSchema(schema schemaStyleGuideConfig) (config policy.StyleGuideConfig) {
	return policy.StyleGuideConfig{
		Path:                schema.Path,
		RequirementIDScheme: requirementid.Scheme(schema.RequirementIDScheme),
	}
}
