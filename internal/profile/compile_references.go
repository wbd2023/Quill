package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func validateConfigRef(
	binding policy.RuleBinding,
	builtin contract.RuleDefinition,
) (err error) {
	if len(builtin.RequiredConfigRefs) == 0 {
		if binding.ConfigRef != "" {
			return fmt.Errorf(
				"rule %q has unexpected config_ref %q",
				binding.RuleID,
				binding.ConfigRef,
			)
		}

		return nil
	}

	for _, configRef := range builtin.RequiredConfigRefs {
		if binding.ConfigRef == configRef {
			return nil
		}
	}

	return fmt.Errorf(
		"rule %q must use config_ref %q",
		binding.RuleID,
		builtin.RequiredConfigRefs[0],
	)
}

func validatePathClasses(
	config policy.Config,
	binding policy.RuleBinding,
) (err error) {
	seen := make(map[string]bool, len(binding.PathClasses))
	for _, className := range binding.PathClasses {
		if className == "" {
			return fmt.Errorf("rule %q has an empty path class", binding.RuleID)
		}

		if seen[className] {
			return fmt.Errorf("rule %q duplicates path class %q", binding.RuleID, className)
		}

		seen[className] = true
		if len(config.Paths.Patterns(className)) == 0 {
			return fmt.Errorf("rule %q references unknown path class %q", binding.RuleID, className)
		}
	}

	return nil
}
