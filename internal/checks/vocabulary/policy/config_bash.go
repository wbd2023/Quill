package policy

func decodeBashConfig(section map[string]any) (config BashConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.vocabulary.bash",
		"forbidden_variable_names",
		"preferred_variable_name",
	); err != nil {
		return BashConfig{}, err
	}

	config.ForbiddenVariableNames, err = stringList(
		section,
		"forbidden_variable_names",
		"packs.vocabulary.bash.forbidden_variable_names",
	)
	if err != nil {
		return BashConfig{}, err
	}

	config.PreferredVariableName, err = stringField(
		section,
		"preferred_variable_name",
		"packs.vocabulary.bash.preferred_variable_name",
	)
	if err != nil {
		return BashConfig{}, err
	}

	return config, nil
}
