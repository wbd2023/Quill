package profile

// Format validates config and returns canonical style profile TOML.
func Format(config Profile) (contents string, err error) {
	if err = Validate(config); err != nil {
		return "", err
	}

	return encodeTOML(config)
}
