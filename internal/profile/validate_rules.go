package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func validateRulePacks(rulePacks policy.RulePackConfig) (err error) {
	if len(rulePacks.Enabled) == 0 {
		return fmt.Errorf("rule_packs.enabled must not be empty")
	}

	return nil
}

func validateRules(
	repository policy.RepositoryConfig,
	rules []policy.RuleBinding,
) (err error) {
	if len(rules) == 0 {
		return fmt.Errorf("rules must not be empty")
	}

	seenRuleIDs := make(map[string]bool, len(rules))
	for _, binding := range rules {
		if binding.RuleID == "" {
			return fmt.Errorf("rule binding has an empty rule_id")
		}

		if seenRuleIDs[binding.RuleID] {
			return fmt.Errorf("duplicate rule binding %q", binding.RuleID)
		}
		seenRuleIDs[binding.RuleID] = true

		switch binding.Level {
		case contract.LevelRequired, contract.LevelRecommendation:
		default:
			return fmt.Errorf("rule %q has invalid level %q", binding.RuleID, binding.Level)
		}

		if !repository.HasScope(binding.Scope) {
			return fmt.Errorf("rule %q references unknown scope %q", binding.RuleID, binding.Scope)
		}

		if len(binding.RequirementIDs) == 0 {
			return fmt.Errorf("rule %q must bind at least one requirement", binding.RuleID)
		}

		seenRequirements := make(map[string]bool, len(binding.RequirementIDs))
		for _, requirementID := range binding.RequirementIDs {
			if requirementID == "" {
				return fmt.Errorf("rule %q has an empty requirement ID", binding.RuleID)
			}

			if seenRequirements[requirementID] {
				return fmt.Errorf(
					"rule %q duplicates requirement %q",
					binding.RuleID,
					requirementID,
				)
			}

			seenRequirements[requirementID] = true
		}
	}

	return nil
}
