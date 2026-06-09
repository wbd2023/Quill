package project

func decodeMakefileVariable(section map[string]any) (variable MakefileVariable, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.project.commands.required_variables",
		"name",
		"value",
	); err != nil {
		return MakefileVariable{}, err
	}

	variable.Name, err = stringField(
		section,
		"name",
		"packs.project.commands.required_variables.name",
	)
	if err != nil {
		return MakefileVariable{}, err
	}

	variable.Value, err = stringField(
		section,
		"value",
		"packs.project.commands.required_variables.value",
	)
	if err != nil {
		return MakefileVariable{}, err
	}

	return variable, nil
}

func decodeMakefileTarget(section map[string]any) (target MakefileTarget, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.project.commands.required_targets",
		"name",
		"recipe_line",
	); err != nil {
		return MakefileTarget{}, err
	}

	target.Name, err = stringField(
		section,
		"name",
		"packs.project.commands.required_targets.name",
	)
	if err != nil {
		return MakefileTarget{}, err
	}

	target.RecipeLine, err = stringField(
		section,
		"recipe_line",
		"packs.project.commands.required_targets.recipe_line",
	)
	if err != nil {
		return MakefileTarget{}, err
	}

	return target, nil
}

func encodeMakefileVariables(variables []MakefileVariable) (tables []map[string]any) {
	tables = make([]map[string]any, 0, len(variables))
	for _, variable := range variables {
		tables = append(tables, map[string]any{
			"name":  variable.Name,
			"value": variable.Value,
		})
	}

	return tables
}

func encodeMakefileTargets(targets []MakefileTarget) (tables []map[string]any) {
	tables = make([]map[string]any, 0, len(targets))
	for _, target := range targets {
		tables = append(tables, map[string]any{
			"name":        target.Name,
			"recipe_line": target.RecipeLine,
		})
	}

	return tables
}
