package profile

import (
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/style"
)

// EffectiveProfile is the resolved Profile and executable rule/tool configuration.
type EffectiveProfile struct {
	Profile   policy.Config
	Effective style.Plan
}

// Compile validates config, applies Pack defaults, and builds an Effective Profile.
func Compile(
	config policy.Config,
	registry pack.Registry,
) (compiled EffectiveProfile, err error) {
	if err := Validate(config); err != nil {
		return EffectiveProfile{}, err
	}

	config, err = effective.ResolvePacks(config, registry.Packs())
	if err != nil {
		return EffectiveProfile{}, err
	}

	if err := Validate(config); err != nil {
		return EffectiveProfile{}, err
	}

	compiled.Effective, err = effective.Compile(config, registry.Definitions())
	if err != nil {
		return EffectiveProfile{}, err
	}

	compiled.Profile = config
	return compiled, nil
}
