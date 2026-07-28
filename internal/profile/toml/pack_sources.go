package toml

import (
	"github.com/wbd2023/quill/internal/policy"
)

// schemaPackSource is the TOML surface for one [[pack_sources]] entry.
type schemaPackSource struct {
	Path string `toml:"path"`
}

func decodePackSources(sources []schemaPackSource) (packSources []policy.PackSource) {
	if len(sources) == 0 {
		return nil
	}

	packSources = make([]policy.PackSource, len(sources))
	for index, source := range sources {
		packSources[index] = policy.PackSource{Path: source.Path}
	}
	return packSources
}

func encodePackSources(sources []policy.PackSource) (schema []schemaPackSource) {
	if len(sources) == 0 {
		return nil
	}

	schema = make([]schemaPackSource, len(sources))
	for index, source := range sources {
		schema[index] = schemaPackSource{Path: source.Path}
	}
	return schema
}
