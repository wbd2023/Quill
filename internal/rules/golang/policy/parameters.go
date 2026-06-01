package policy

// ParameterConfig defines Go parameter classification policy.
type ParameterConfig struct {
	SecretNames []string
}

func decodeParameterConfig(
	section map[string]any,
) (config ParameterConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.go.parameters",
		"secret_names",
	); err != nil {
		return ParameterConfig{}, err
	}

	config.SecretNames, err = stringList(
		section,
		"secret_names",
		"packs.go.parameters.secret_names",
	)
	if err != nil {
		return ParameterConfig{}, err
	}

	return config, nil
}

func validateParameters(config ParameterConfig) (err error) {
	return validateList("packs.go.parameters.secret_names", config.SecretNames)
}
