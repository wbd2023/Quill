package projectpolicy

func decodeMakeConfig(section map[string]any) (config MakeConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.project.commands",
		"runner",
		"path",
		"required_variables",
		"required_targets",
	); err != nil {
		return MakeConfig{}, err
	}

	config.Path, err = stringField(section, "path", "packs.project.commands.path")
	if err != nil {
		return MakeConfig{}, err
	}

	variables, err := tableList(
		section,
		"required_variables",
		"packs.project.commands.required_variables",
	)
	if err != nil {
		return MakeConfig{}, err
	}

	config.RequiredVariables = make([]MakefileVariable, 0, len(variables))
	for _, variable := range variables {
		required, err := decodeMakefileVariable(variable)
		if err != nil {
			return MakeConfig{}, err
		}

		config.RequiredVariables = append(config.RequiredVariables, required)
	}

	targets, err := tableList(
		section,
		"required_targets",
		"packs.project.commands.required_targets",
	)
	if err != nil {
		return MakeConfig{}, err
	}

	config.RequiredTargets = make([]MakefileTarget, 0, len(targets))
	for _, target := range targets {
		required, err := decodeMakefileTarget(target)
		if err != nil {
			return MakeConfig{}, err
		}

		config.RequiredTargets = append(config.RequiredTargets, required)
	}

	return config, nil
}
