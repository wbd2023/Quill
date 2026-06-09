package vocabulary

import "fmt"

// ValidateConfig validates project vocabulary policy.
func ValidateConfig(config Config) (err error) {
	goConfig := config.Go
	bashConfig := config.Bash

	if err = validateList(
		"packs.vocabulary.go.forbidden_type_suffixes",
		goConfig.ForbiddenTypeSuffixes,
	); err != nil {
		return err
	}

	if err = validateList(
		"packs.vocabulary.go.forbidden_identifier_suffixes",
		goConfig.ForbiddenIdentifierSuffixes,
	); err != nil {
		return err
	}

	if err = validateList(
		"packs.vocabulary.bash.forbidden_variable_names",
		bashConfig.ForbiddenVariableNames,
	); err != nil {
		return err
	}

	if len(goConfig.ForbiddenTypeSuffixes) > 0 && goConfig.PreferredTypeSuffix == "" {
		return fmt.Errorf("packs.vocabulary.go.preferred_type_suffix must not be empty")
	}

	if len(goConfig.ForbiddenIdentifierSuffixes) > 0 &&
		goConfig.PreferredIdentifierSuffix == "" {
		return fmt.Errorf("packs.vocabulary.go.preferred_identifier_suffix must not be empty")
	}

	if len(bashConfig.ForbiddenVariableNames) > 0 &&
		bashConfig.PreferredVariableName == "" {
		return fmt.Errorf("packs.vocabulary.bash.preferred_variable_name must not be empty")
	}

	return nil
}
