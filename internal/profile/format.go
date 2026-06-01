package profile

import (
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/toml"
)

// Format validates config and returns canonical style profile TOML.
func Format(config policy.Config) (contents string, err error) {
	if err = Validate(config); err != nil {
		return "", err
	}

	return toml.Encode(config)
}
