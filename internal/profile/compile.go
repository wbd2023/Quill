package profile

import (
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

// EffectiveProfile is the resolved Profile and executable rule/tool configuration.
type EffectiveProfile struct {
	Profile   policy.Config
	Effective style.Plan
}

// Compile validates config and builds an executable plan from definitions.
func Compile(
	config policy.Config,
	definitions style.Definitions,
) (compiled EffectiveProfile, err error) {
	if err := Validate(config); err != nil {
		return EffectiveProfile{}, err
	}

	compiled.Effective, err = compilePlan(config, definitions)
	if err != nil {
		return EffectiveProfile{}, err
	}

	compiled.Profile = config
	return compiled, nil
}
