package toml

import (
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

type schemaRepositoryConfig struct {
	RootMarkers         []string            `toml:"root_markers"`
	DefaultScope        string              `toml:"default_scope"`
	ScopeRoots          map[string][]string `toml:"scope_roots"`
	ExcludedDirectories []string            `toml:"excluded_directories"`
	GeneratedMarker     string              `toml:"generated_marker"`
}

func decodeRepository(schema schemaRepositoryConfig) (repository policy.RepositoryConfig) {
	return policy.RepositoryConfig{
		RootMarkers:         append([]string{}, schema.RootMarkers...),
		DefaultScope:        style.Scope(schema.DefaultScope),
		ScopeRoots:          decodeScopeMap(schema.ScopeRoots),
		ExcludedDirectories: append([]string{}, schema.ExcludedDirectories...),
		GeneratedMarker:     schema.GeneratedMarker,
	}
}

func encodeRepository(repository policy.RepositoryConfig) (schema schemaRepositoryConfig) {
	return schemaRepositoryConfig{
		RootMarkers:         append([]string{}, repository.RootMarkers...),
		DefaultScope:        string(repository.DefaultScope),
		ScopeRoots:          encodeScopeMap(repository.ScopeRoots),
		ExcludedDirectories: append([]string{}, repository.ExcludedDirectories...),
		GeneratedMarker:     repository.GeneratedMarker,
	}
}
