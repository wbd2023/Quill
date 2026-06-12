package projectpolicy

import (
	"fmt"
	"strings"
)

// ValidateConfig validates Project Pack Policy.
func ValidateConfig(config Config) (err error) {
	return validateCommands(config.Commands)
}

func validateCommands(commands CommandsConfig) (err error) {
	switch {
	case blank(string(commands.Runner)):
		return fmt.Errorf("packs.project.commands.runner must not be empty")
	case commands.Runner == CommandsRunnerMake:
		return validateMakeConfig(commands.Make)
	default:
		return fmt.Errorf("unsupported packs.project.commands.runner %q", commands.Runner)
	}
}

func validateMakeConfig(config MakeConfig) (err error) {
	if blank(config.Path) {
		return fmt.Errorf("packs.project.commands.path must not be empty")
	}

	if len(config.RequiredTargets) == 0 {
		return fmt.Errorf(
			"packs.project.commands.required_targets must not be empty",
		)
	}

	seenVariables := make(map[string]bool, len(config.RequiredVariables))
	for _, variable := range config.RequiredVariables {
		if blank(variable.Name) {
			return fmt.Errorf(
				"packs.project.commands.required_variables contains an empty name",
			)
		}

		if seenVariables[variable.Name] {
			return fmt.Errorf(
				"packs.project.commands.required_variables contains duplicate name %q",
				variable.Name,
			)
		}

		seenVariables[variable.Name] = true
	}

	seenTargets := make(map[string]bool, len(config.RequiredTargets))
	for _, target := range config.RequiredTargets {
		if blank(target.Name) {
			return fmt.Errorf(
				"packs.project.commands.required_targets contains an empty name",
			)
		}

		if seenTargets[target.Name] {
			return fmt.Errorf(
				"packs.project.commands.required_targets contains duplicate name %q",
				target.Name,
			)
		}

		seenTargets[target.Name] = true
		if blank(target.RecipeLine) {
			targetField := "packs.project.commands.required_targets." +
				target.Name
			return fmt.Errorf("%s recipe_line must not be empty", targetField)
		}
	}

	return nil
}

func blank(value string) (blank bool) {
	return strings.TrimSpace(value) == ""
}
