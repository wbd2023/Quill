package vocabularypolicy

func decodeGoConfig(section map[string]any) (config GoConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.vocabulary.go",
		"type_suffixes",
		"identifier_suffixes",
	); err != nil {
		return GoConfig{}, err
	}

	config.TypeSuffixes, err = stringListMap(
		section,
		"type_suffixes",
		"packs.vocabulary.go.type_suffixes",
	)
	if err != nil {
		return GoConfig{}, err
	}

	config.IdentifierSuffixes, err = stringListMap(
		section,
		"identifier_suffixes",
		"packs.vocabulary.go.identifier_suffixes",
	)
	if err != nil {
		return GoConfig{}, err
	}

	return config, nil
}
