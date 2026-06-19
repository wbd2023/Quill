package validation

import (
	"fmt"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

func validateRules(
	repository policy.RepositoryConfig,
	scheme style.IDScheme,
	rules []policy.RuleBinding,
) (err error) {
	if len(rules) == 0 {
		return fmt.Errorf("rules must not be empty")
	}

	seenRules := make(map[string]bool, len(rules))
	for _, binding := range rules {
		if isBlank(binding.RuleID) {
			return fmt.Errorf("rule binding has an empty id")
		}

		if seenRules[binding.RuleID] {
			return fmt.Errorf("duplicate rule binding %q", binding.RuleID)
		}
		seenRules[binding.RuleID] = true

		switch binding.Enforcement {
		case style.EnforcementRequired, style.EnforcementRecommendation:

		default:
			return fmt.Errorf(
				"rule %q has invalid enforcement %q",
				binding.RuleID,
				binding.Enforcement,
			)
		}

		if !repository.HasScope(binding.Scope) {
			return fmt.Errorf("rule %q references unknown scope %q", binding.RuleID, binding.Scope)
		}

		if len(binding.RequirementIDs) == 0 {
			return fmt.Errorf("rule %q must bind at least one requirement", binding.RuleID)
		}

		seenRequirements := make(map[string]bool, len(binding.RequirementIDs))
		for _, id := range binding.RequirementIDs {
			if isBlank(id) {
				return fmt.Errorf("rule %q has an empty requirement id", binding.RuleID)
			}

			if seenRequirements[id] {
				return fmt.Errorf(
					"rule %q duplicates requirement %q",
					binding.RuleID,
					id,
				)
			}

			if _, err = style.ParseRequirementID(id, scheme); err != nil {
				return fmt.Errorf(
					"rule %q has invalid requirement id %q: %w",
					binding.RuleID,
					id,
					err,
				)
			}

			seenRequirements[id] = true
		}
	}

	return nil
}
