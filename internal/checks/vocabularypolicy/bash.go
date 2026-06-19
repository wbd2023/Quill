package vocabularypolicy

func decodeBashConfig(section map[string]any) (config BashConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.vocabulary.bash",
		"variable_names",
	); err != nil {
		return BashConfig{}, err
	}

	config.VariableNames, err = stringListMap(
		section,
		"variable_names",
		"packs.vocabulary.bash.variable_names",
	)
	if err != nil {
		return BashConfig{}, err
	}

	return config, nil
}
