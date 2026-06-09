package project

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

// DecodeConfig decodes the project pack config subtree.
func DecodeConfig(pack policy.PackConfig) (config Config, err error) {
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

// ValidatePackConfig validates the raw project pack config subtree.
func ValidatePackConfig(pack policy.PackConfig) (err error) {
	_, err = DecodeConfig(pack)
	return err
}

// EncodeConfig encodes config as a raw project pack config subtree.
func EncodeConfig(config Config) (pack policy.PackConfig) {
	makeConfig := config.Commands.Make

	return policy.PackConfig{
		"commands": map[string]any{
			"runner":             string(config.Commands.Runner),
			"path":               makeConfig.Path,
			"required_variables": encodeMakefileVariables(makeConfig.RequiredVariables),
			"required_targets":   encodeMakefileTargets(makeConfig.RequiredTargets),
		},
	}
}
