package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func validateConfigReference(
	binding policy.RuleBinding,
	definition contract.RuleDefinition,
) (err error) {
	if len(definition.RequiredConfigReferences) == 0 {
		if binding.ConfigReference != "" {
			return fmt.Errorf(
				"rule %q has unexpected config_reference %q",
				binding.RuleID,
				binding.ConfigReference,
			)
		}

		return nil
	}

	for _, configReference := range definition.RequiredConfigReferences {
		if binding.ConfigReference == configReference {
			return nil
		}
	}

	return fmt.Errorf(
		"rule %q must use config_reference %q",
		binding.RuleID,
		definition.RequiredConfigReferences[0],
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
		if len(config.Paths.LookupPatterns(className)) == 0 {
			return fmt.Errorf("rule %q references unknown path class %q", binding.RuleID, className)
		}
	}

	return nil
}
