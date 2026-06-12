package gopolicy

import "fmt"

func validateArchitecture(architecture ArchitectureConfig) (err error) {
	if len(architecture.Layers) == 0 {
		return fmt.Errorf("packs.go.architecture.layers must not be empty")
	}

	knownLayers := make(map[string]bool, len(architecture.Layers))
	for _, layer := range architecture.Layers {
		if blank(layer.Name) {
			return fmt.Errorf("packs.go.architecture.layers contains an empty name")
		}

		if knownLayers[layer.Name] {
			return fmt.Errorf(
				"packs.go.architecture.layers contains duplicate layer %q",
				layer.Name,
			)
		}

		knownLayers[layer.Name] = true

		if len(layer.PackageRoots) == 0 {
			return fmt.Errorf(
				"packs.go.architecture.layers.%s.package_roots must not be empty",
				layer.Name,
			)
		}

		if err = validateList(
			"packs.go.architecture.layers."+layer.Name+".package_roots",
			layer.PackageRoots,
		); err != nil {
			return err
		}
	}

	for _, layer := range architecture.Layers {
		if err = validateList(
			"packs.go.architecture.layers."+layer.Name+".may_import",
			layer.AllowedLayers,
		); err != nil {
			return err
		}

		for _, allowedLayer := range layer.AllowedLayers {
			if knownLayers[allowedLayer] {
				continue
			}

			return fmt.Errorf(
				"packs.go.architecture.layers.%s references unknown layer %q",
				layer.Name,
				allowedLayer,
			)
		}
	}

	return nil
}
