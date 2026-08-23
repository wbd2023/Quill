package policy

import (
	"fmt"

	"github.com/wbd2023/quill/internal/profile"
)

/* ------------------------------------------ Decoding ------------------------------------------ */

// DecodeConfig decodes the Go Pack Policy subtree.
func DecodeConfig(policy profile.PackPolicy) (config Config, err error) {
	if policy == nil {
		return Config{}, fmt.Errorf("packs.go must be configured")
	}

	if err = rejectUnknownFields(
		policy,
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
		policy,
		"local_import_prefixes",
		"packs.go.local_import_prefixes",
	)
	if err != nil {
		return Config{}, err
	}

	parameters, err := configSection(policy, "parameters", "packs.go.parameters")
	if err != nil {
		return Config{}, err
	}

	config.Parameters, err = decodeParameterConfig(parameters)
	if err != nil {
		return Config{}, err
	}

	constructors, err := configSection(policy, "constructors", "packs.go.constructors")
	if err != nil {
		return Config{}, err
	}

	config.Constructors, err = decodeConstructorConfig(constructors)
	if err != nil {
		return Config{}, err
	}

	domainValues, err := configSection(
		policy,
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

	architecture, err := configSection(policy, "architecture", "packs.go.architecture")
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

// ValidatePackPolicy validates the raw Go Pack Policy subtree.
func ValidatePackPolicy(policy profile.PackPolicy) (err error) {
	_, err = DecodeConfig(policy)
	return err
}

/* ------------------------------------------ Encoding ------------------------------------------ */

// EncodeConfig encodes config as a raw Go Pack Policy subtree.
func EncodeConfig(config Config) (policy profile.PackPolicy) {
	return profile.PackPolicy{
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
