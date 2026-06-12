package gopolicy

import (
	"fmt"

	corepolicy "ciphera/tools/internal/policy"
)

/* ------------------------------------------ Decoding ------------------------------------------ */

// DecodeConfig decodes the Go pack config subtree.
func DecodeConfig(pack corepolicy.PackConfig) (config Config, err error) {
	if pack == nil {
		return Config{}, fmt.Errorf("packs.go must be configured")
	}

	if err = rejectUnknownFields(
		pack,
		"packs.go",
		"local_import_prefixes",
		"parameters",
		"constructors",
		"domain_values",
		"architecture",
	); err != nil {
		return Config{}, err
	}

	config.LocalImportPrefixes, err = stringList(
		pack,
		"local_import_prefixes",
		"packs.go.local_import_prefixes",
	)
	if err != nil {
		return Config{}, err
	}

	parameters, err := configSection(pack, "parameters", "packs.go.parameters")
	if err != nil {
		return Config{}, err
	}

	config.Parameters, err = decodeParameterConfig(parameters)
	if err != nil {
		return Config{}, err
	}

	constructors, err := configSection(pack, "constructors", "packs.go.constructors")
	if err != nil {
		return Config{}, err
	}

	config.Constructors, err = decodeConstructorConfig(constructors)
	if err != nil {
		return Config{}, err
	}

	domainValues, err := configSection(
		pack,
		"domain_values",
		"packs.go.domain_values",
	)
	if err != nil {
		return Config{}, err
	}

	config.DomainValues.RequiredConstructors, err = stringListMap(
		domainValues,
		"required_constructors",
		"packs.go.domain_values.required_constructors",
	)
	if err != nil {
		return Config{}, err
	}

	architecture, err := configSection(pack, "architecture", "packs.go.architecture")
	if err != nil {
		return Config{}, err
	}

	config.Architecture, err = decodeArchitectureConfig(architecture)
	if err != nil {
		return Config{}, err
	}

	return config, ValidateConfig(config)
}

/* ----------------------------------------- Validation ----------------------------------------- */

// ValidatePackConfig validates the raw Go pack config subtree.
func ValidatePackConfig(pack corepolicy.PackConfig) (err error) {
	_, err = DecodeConfig(pack)
	return err
}

/* ------------------------------------------ Encoding ------------------------------------------ */

// EncodeConfig encodes config as a raw Go pack config subtree.
func EncodeConfig(config Config) (pack corepolicy.PackConfig) {
	return corepolicy.PackConfig{
		"local_import_prefixes": cloneStrings(config.LocalImportPrefixes),
		"parameters": map[string]any{
			"secret_names": cloneStrings(config.Parameters.SecretNames),
		},
		"constructors": map[string]any{
			"parameter_order": encodeParameterGroups(config.Constructors.ParameterOrder),
		},
		"domain_values": map[string]any{
			"required_constructors": encodeStringListMap(
				config.DomainValues.RequiredConstructors,
			),
		},
		"architecture": map[string]any{
			"layers": encodeArchitectureLayers(config.Architecture.Layers),
		},
	}
}
