package profile

// schemaPackSource is the TOML surface for one [[pack_sources]] entry.
type schemaPackSource struct {
	Path string `toml:"path"`
}

func decodePackSources(sources []schemaPackSource) (packSources []PackSource) {
	if len(sources) == 0 {
		return nil
	}

	packSources = make([]PackSource, len(sources))
	for index, source := range sources {
		packSources[index] = PackSource(source)
	}
	return packSources
}

func encodePackSources(sources []PackSource) (schema []schemaPackSource) {
	if len(sources) == 0 {
		return nil
	}

	schema = make([]schemaPackSource, len(sources))
	for index, source := range sources {
		schema[index] = schemaPackSource(source)
	}
	return schema
}
