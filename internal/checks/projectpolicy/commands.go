package projectpolicy

func decodeCommands(
	section map[string]any,
) (commands CommandsConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.project.commands",
		"runner",
		"path",
		"required_variables",
		"required_targets",
	); err != nil {
		return CommandsConfig{}, err
	}

	runner, err := stringField(section, "runner", "packs.project.commands.runner")
	if err != nil {
		return CommandsConfig{}, err
	}

	commands.Runner = CommandsRunner(runner)
	commands.Make, err = decodeMakeConfig(section)
	if err != nil {
		return CommandsConfig{}, err
	}

	return commands, nil
}
