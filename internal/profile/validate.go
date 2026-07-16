package profile

import "ciphera/tools/internal/policy"

// Validate checks config for supported schema version and internal consistency.
func Validate(config policy.Config) (err error) {
	return check(config)
}
