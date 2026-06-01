package profile

import (
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/validation"
)

// Validate checks config for supported schema version and internal consistency.
func Validate(config policy.Config) (err error) {
	return validation.Check(config)
}
