package policy

import (
	"fmt"

	"github.com/wbd2023/quill/internal/profile"
)

// DecodeConfig decodes the Project Pack Policy subtree.
func DecodeConfig(policy profile.PackPolicy) (config Config, err error) {
	if policy == nil {
		return Config{}, fmt.Errorf("packs.project must be configured")
	}

	if err = rejectUnknownFields(policy, "packs.project", "commands"); err != nil {
		return Config{}, err
	}

	section, err := configSection(
		policy,
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

// Validate validates the raw Project Pack Policy subtree.
func Validate(policy profile.PackPolicy) (err error) {
	_, err = DecodeConfig(policy)
	return err
}

// EncodeConfig encodes config as a raw Project Pack Policy subtree.
func EncodeConfig(config Config) (policy profile.PackPolicy) {
	makeConfig := config.Commands.Make
	return profile.PackPolicy{
		"commands": map[string]any{
			"runner":             string(config.Commands.Runner),
			"path":               makeConfig.Path,
			"required_variables": encodeMakefileVariables(makeConfig.RequiredVariables),
			"required_targets":   encodeMakefileTargets(makeConfig.RequiredTargets),
		},
	}
}
