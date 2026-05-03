package profile

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/policy"
)

/* ------------------------------------------ Go Policy ----------------------------------------- */

func validateGo(goConfig policy.GoConfig, language policy.LanguageConfig) (err error) {
	if !hasLanguageBackend(language, "go") {
		return nil
	}

	if len(goConfig.LocalImportPrefixes) == 0 {
		return fmt.Errorf(
			"go.local_import_prefixes must not be empty when Go backends are configured",
		)
	}

	for _, prefix := range goConfig.LocalImportPrefixes {
		if strings.TrimSpace(prefix) == "" {
			return fmt.Errorf("go.local_import_prefixes contains an empty prefix")
		}
	}

	if err = validateGoDomainIdentifierConstructors(
		goConfig.DomainIdentifierConstructors,
	); err != nil {
		return err
	}

	return validateGoArchitecture(goConfig.Architecture)
}

/* ------------------------------------- Domain Identifiers ------------------------------------- */

func validateGoDomainIdentifierConstructors(
	constructors policy.GoDomainIdentifierConstructors,
) (err error) {
	for typeName, typeConstructors := range constructors {
		if typeName == "" {
			return fmt.Errorf("go.domain_identifiers has an empty type name")
		}

		if len(typeConstructors) == 0 {
			return fmt.Errorf(
				"go.domain_identifiers.%s must define at least one constructor",
				typeName,
			)
		}

		for _, constructor := range typeConstructors {
			if constructor == "" {
				return fmt.Errorf(
					"go.domain_identifiers.%s contains an empty constructor",
					typeName,
				)
			}
		}
	}

	return nil
}

/* ---------------------------------------- Architecture ---------------------------------------- */

func validateGoArchitecture(architecture policy.GoArchitectureConfig) (err error) {
	if len(architecture.Layers) == 0 {
		return fmt.Errorf("go.architecture.layers must not be empty")
	}

	knownLayers := make(map[string]bool, len(architecture.Layers))
	for _, layer := range architecture.Layers {
		if layer.Name == "" {
			return fmt.Errorf("go architecture layer name must not be empty")
		}

		if knownLayers[layer.Name] {
			return fmt.Errorf("duplicate go architecture layer %q", layer.Name)
		}

		knownLayers[layer.Name] = true

		if len(layer.PackageRoots) == 0 {
			return fmt.Errorf("go architecture layer %q must define package_roots", layer.Name)
		}
	}

	for _, layer := range architecture.Layers {
		for _, allowedLayer := range layer.AllowedLayers {
			if knownLayers[allowedLayer] {
				continue
			}

			return fmt.Errorf(
				"go architecture layer %q references unknown layer %q",
				layer.Name,
				allowedLayer,
			)
		}
	}

	return nil
}

/* -------------------------------------- Language Backends ------------------------------------- */

func hasLanguageBackend(language policy.LanguageConfig, target string) (found bool) {
	for _, backend := range language.Backends {
		if backend.Language == target {
			return true
		}
	}

	return false
}
