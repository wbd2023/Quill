package policy

// ArchitectureConfig defines package import layering policy.
type ArchitectureConfig struct {
	Layers []ArchitectureLayer
}

// ArchitectureLayer defines one package import layer.
type ArchitectureLayer struct {
	Name          string
	PackageRoots  []string
	AllowedLayers []string
}

func decodeArchitectureConfig(
	section map[string]any,
) (config ArchitectureConfig, err error) {
	if err = rejectUnknownFields(section, "packs.go.architecture", "layers"); err != nil {
		return ArchitectureConfig{}, err
	}

	layers, err := tableList(section, "layers", "packs.go.architecture.layers")
	if err != nil {
		return ArchitectureConfig{}, err
	}

	config.Layers = make([]ArchitectureLayer, 0, len(layers))
	for _, layer := range layers {
		architectureLayer, err := decodeArchitectureLayer(layer)
		if err != nil {
			return ArchitectureConfig{}, err
		}

		config.Layers = append(config.Layers, architectureLayer)
	}

	return config, nil
}

func decodeArchitectureLayer(section map[string]any) (layer ArchitectureLayer, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.go.architecture.layers",
		"name",
		"package_roots",
		"may_import",
	); err != nil {
		return ArchitectureLayer{}, err
	}

	layer.Name, err = stringField(section, "name", "packs.go.architecture.layers.name")
	if err != nil {
		return ArchitectureLayer{}, err
	}

	layer.PackageRoots, err = stringList(
		section,
		"package_roots",
		"packs.go.architecture.layers.package_roots",
	)
	if err != nil {
		return ArchitectureLayer{}, err
	}

	layer.AllowedLayers, err = stringList(
		section,
		"may_import",
		"packs.go.architecture.layers.may_import",
	)
	if err != nil {
		return ArchitectureLayer{}, err
	}

	return layer, nil
}

func encodeArchitectureLayers(layers []ArchitectureLayer) (tables []map[string]any) {
	tables = make([]map[string]any, 0, len(layers))
	for _, layer := range layers {
		tables = append(tables, map[string]any{
			"name":          layer.Name,
			"package_roots": cloneStrings(layer.PackageRoots),
			"may_import":    cloneStrings(layer.AllowedLayers),
		})
	}

	return tables
}
