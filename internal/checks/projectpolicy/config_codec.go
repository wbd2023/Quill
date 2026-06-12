package projectpolicy

import (
	"fmt"

	corepolicy "ciphera/tools/internal/policy"
)

// DecodeConfig decodes the Project Pack Policy subtree.
func DecodeConfig(pack corepolicy.PackConfig) (config Config, err error) {
	if pack == nil {
		return Config{}, fmt.Errorf("packs.project must be configured")
	}

	if err = rejectUnknownFields(pack, "packs.project", "commands"); err != nil {
		return Config{}, err
	}

	section, err := configSection(
		pack,
		"commands",
		"packs.project.commands",
	)
	if err != nil {
		return Config{}, err
	}

	config.Commands, err = decodeCommands(section)
	if err != nil {
		return Config{}, err
	}

	return config, ValidateConfig(config)
}

// ValidatePackConfig validates the raw Project Pack Policy subtree.
func ValidatePackConfig(pack corepolicy.PackConfig) (err error) {
	_, err = DecodeConfig(pack)
	return err
}

// EncodeConfig encodes config as a raw Project Pack Policy subtree.
func EncodeConfig(config Config) (pack corepolicy.PackConfig) {
	makeConfig := config.Commands.Make

	return corepolicy.PackConfig{
		"commands": map[string]any{
			"runner":             string(config.Commands.Runner),
			"path":               makeConfig.Path,
			"required_variables": encodeMakefileVariables(makeConfig.RequiredVariables),
			"required_targets":   encodeMakefileTargets(makeConfig.RequiredTargets),
		},
	}
}
