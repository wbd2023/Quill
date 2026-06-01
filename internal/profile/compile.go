package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/effective"
)

// Compile validates config and resolves it against available rule and tool definitions.
func Compile(
	config policy.Config,
	definitions contract.Definitions,
) (compiled contract.EffectiveConfig, err error) {
	if err := Validate(config); err != nil {
		return contract.EffectiveConfig{}, err
	}

	return effective.Compile(config, definitions)
}
