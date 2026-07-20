package profile

import (
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/profile/toml"
)

// Format validates config and returns canonical style profile TOML.
func Format(config policy.Config) (contents string, err error) {
	if err = Validate(config); err != nil {
		return "", err
	}

	return toml.Encode(config)
}
