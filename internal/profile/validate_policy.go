package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

/* ------------------------------------------ Language ------------------------------------------ */

func validateLanguage(
	repository policy.RepositoryConfig,
	language policy.LanguageConfig,
) (err error) {
	seenLanguageBackends := make(map[string]bool, len(language.Backends))
	for _, backend := range language.Backends {
		if backend.Name == "" {
			return fmt.Errorf("language backend name must not be empty")
		}

		if seenLanguageBackends[backend.Name] {
			return fmt.Errorf("duplicate language backend %q", backend.Name)
		}

		seenLanguageBackends[backend.Name] = true

		if backend.Language == "" {
			return fmt.Errorf("language backend %q must define language", backend.Name)
		}

		if !repository.ScopeExists(backend.Scope) {
			return fmt.Errorf(
				"language backend %q references unknown scope %q",
				backend.Name,
				backend.Scope,
			)
		}
	}

	return nil
}

/* -------------------------------------------- Tools ------------------------------------------- */

func validateTools(tools []policy.ToolPin) (err error) {
	seenToolPins := make(map[string]bool, len(tools))
	for _, tool := range tools {
		if tool.ID == "" {
			return fmt.Errorf("tool pin has an empty id")
		}

		if seenToolPins[tool.ID] {
			return fmt.Errorf("duplicate tool pin %q", tool.ID)
		}

		seenToolPins[tool.ID] = true
		if tool.Version == "" {
			return fmt.Errorf("tool pin %q must define version", tool.ID)
		}
	}

	return nil
}

/* ------------------------------------------- Naming ------------------------------------------- */

func validateNaming(naming policy.NamingConfig) (err error) {
	return validateGoDomainIdentifiers(naming.GoDomainIdentifiers)
}

func validateGoDomainIdentifiers(identifiers policy.GoDomainIdentifierConfig) (err error) {
	for typeName, constructors := range identifiers {
		if typeName == "" {
			return fmt.Errorf("naming.go_domain_identifiers has an empty type name")
		}

		if len(constructors) == 0 {
			return fmt.Errorf(
				"naming.go_domain_identifiers.%s must define at least one constructor",
				typeName,
			)
		}

		for _, constructor := range constructors {
			if constructor == "" {
				return fmt.Errorf(
					"naming.go_domain_identifiers.%s contains an empty constructor",
					typeName,
				)
			}
		}
	}

	return nil
}

/* ---------------------------------------- Control Plane --------------------------------------- */

func validateControlPlane(controlPlane policy.ControlPlaneConfig) (err error) {
	if controlPlane.QualityFile == "" {
		return fmt.Errorf("control_plane.quality_file must not be empty")
	}

	return nil
}

/* ---------------------------------------- Architecture ---------------------------------------- */

func validateArchitecture(architecture policy.ArchitectureConfig) (err error) {
	if len(architecture.Layers) == 0 {
		return fmt.Errorf("architecture.layers must not be empty")
	}

	layerNames := make(map[string]bool, len(architecture.Layers))
	for _, layer := range architecture.Layers {
		if layer.Name == "" {
			return fmt.Errorf("architecture layer name must not be empty")
		}

		if layerNames[layer.Name] {
			return fmt.Errorf("duplicate architecture layer %q", layer.Name)
		}

		layerNames[layer.Name] = true

		if len(layer.PackageRoots) == 0 {
			return fmt.Errorf("architecture layer %q must define package_roots", layer.Name)
		}
	}

	for _, layer := range architecture.Layers {
		for _, allowed := range layer.MayImport {
			if layerNames[allowed] {
				continue
			}

			return fmt.Errorf(
				"architecture layer %q references unknown layer %q",
				layer.Name,
				allowed,
			)
		}
	}

	return nil
}
