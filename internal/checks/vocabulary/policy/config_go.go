package policy

func decodeGoConfig(section map[string]any) (config GoConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.vocabulary.go",
		"forbidden_type_suffixes",
		"preferred_type_suffix",
		"forbidden_identifier_suffixes",
		"preferred_identifier_suffix",
	); err != nil {
		return GoConfig{}, err
	}

	config.ForbiddenTypeSuffixes, err = stringList(
		section,
		"forbidden_type_suffixes",
		"packs.vocabulary.go.forbidden_type_suffixes",
	)
	if err != nil {
		return GoConfig{}, err
	}

	config.PreferredTypeSuffix, err = stringField(
		section,
		"preferred_type_suffix",
		"packs.vocabulary.go.preferred_type_suffix",
	)
	if err != nil {
		return GoConfig{}, err
	}

	config.ForbiddenIdentifierSuffixes, err = stringList(
		section,
		"forbidden_identifier_suffixes",
		"packs.vocabulary.go.forbidden_identifier_suffixes",
	)
	if err != nil {
		return GoConfig{}, err
	}

	config.PreferredIdentifierSuffix, err = stringField(
		section,
		"preferred_identifier_suffix",
		"packs.vocabulary.go.preferred_identifier_suffix",
	)
	if err != nil {
		return GoConfig{}, err
	}

	return config, nil
}
