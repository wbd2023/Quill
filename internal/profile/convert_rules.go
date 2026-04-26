package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func rulesFromSchema(schemas []schemaRuleBinding) (rules []policy.RuleBinding) {
	rules = make([]policy.RuleBinding, 0, len(schemas))
	for _, rule := range schemas {
		rules = append(rules, policy.RuleBinding{
			RuleID:         rule.RuleID,
			Level:          rule.Level,
			Scope:          contract.Scope(rule.Scope),
			RequirementIDs: append([]string{}, rule.RequirementIDs...),
			ConfigRef:      rule.ConfigRef,
			Backends:       append([]string{}, rule.Backends...),
			PathClasses:    append([]string{}, rule.PathClasses...),
		})
	}

	return rules
}

func rulesToSchema(rules []policy.RuleBinding) (schemas []schemaRuleBinding) {
	schemas = make([]schemaRuleBinding, 0, len(rules))
	for _, rule := range rules {
		schemas = append(schemas, schemaRuleBinding{
			RuleID:         rule.RuleID,
			Level:          rule.Level,
			Scope:          string(rule.Scope),
			RequirementIDs: append([]string{}, rule.RequirementIDs...),
			ConfigRef:      rule.ConfigRef,
			Backends:       append([]string{}, rule.Backends...),
			PathClasses:    append([]string{}, rule.PathClasses...),
		})
	}

	return schemas
}
