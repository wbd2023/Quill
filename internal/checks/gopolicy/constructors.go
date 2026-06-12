package gopolicy

/* -------------------------------------- Constructor Types ------------------------------------- */

// ConstructorConfig defines Go constructor policy.
type ConstructorConfig struct {
	ParameterOrder []ParameterGroup
}

// ParameterGroup defines one constructor parameter ordering class.
type ParameterGroup struct {
	Name               string
	TypeNameSuffixes   []string
	ParameterNames     []string
	MatchesSecretNames bool
}

/* ------------------------------------------ Decoding ------------------------------------------ */

func decodeConstructorConfig(
	section map[string]any,
) (config ConstructorConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.go.constructors",
		"parameter_order",
	); err != nil {
		return ConstructorConfig{}, err
	}

	groups, err := tableList(
		section,
		"parameter_order",
		"packs.go.constructors.parameter_order",
	)
	if err != nil {
		return ConstructorConfig{}, err
	}

	config.ParameterOrder = make([]ParameterGroup, 0, len(groups))
	for _, group := range groups {
		parameterGroup, err := decodeParameterGroup(group)
		if err != nil {
			return ConstructorConfig{}, err
		}

		config.ParameterOrder = append(config.ParameterOrder, parameterGroup)
	}

	return config, nil
}

func decodeParameterGroup(section map[string]any) (group ParameterGroup, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.go.constructors.parameter_order",
		"name",
		"type_name_suffixes",
		"parameter_names",
		"matches_secret_names",
	); err != nil {
		return ParameterGroup{}, err
	}

	group.Name, err = stringField(section, "name", "packs.go.constructors.parameter_order.name")
	if err != nil {
		return ParameterGroup{}, err
	}

	group.TypeNameSuffixes, err = stringList(
		section,
		"type_name_suffixes",
		"packs.go.constructors.parameter_order.type_name_suffixes",
	)
	if err != nil {
		return ParameterGroup{}, err
	}

	group.ParameterNames, err = stringList(
		section,
		"parameter_names",
		"packs.go.constructors.parameter_order.parameter_names",
	)
	if err != nil {
		return ParameterGroup{}, err
	}

	group.MatchesSecretNames, err = boolField(
		section,
		"matches_secret_names",
		"packs.go.constructors.parameter_order.matches_secret_names",
	)
	if err != nil {
		return ParameterGroup{}, err
	}

	return group, nil
}

/* ------------------------------------------ Encoding ------------------------------------------ */

func encodeParameterGroups(groups []ParameterGroup) (tables []map[string]any) {
	tables = make([]map[string]any, 0, len(groups))
	for _, group := range groups {
		tables = append(tables, map[string]any{
			"name":                 group.Name,
			"type_name_suffixes":   cloneStrings(group.TypeNameSuffixes),
			"parameter_names":      cloneStrings(group.ParameterNames),
			"matches_secret_names": group.MatchesSecretNames,
		})
	}

	return tables
}
