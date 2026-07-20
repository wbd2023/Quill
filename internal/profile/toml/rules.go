package toml

import (
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

type schemaRuleBinding struct {
	ID             string            `toml:"id"`
	Enforcement    style.Enforcement `toml:"enforcement"`
	Scope          string            `toml:"scope"`
	RequirementIDs []string          `toml:"requirement_ids"`
}

func decodeRules(schemas []schemaRuleBinding) (rules []policy.RuleBinding) {
	rules = make([]policy.RuleBinding, 0, len(schemas))
	for _, rule := range schemas {
		rules = append(rules, policy.RuleBinding{
			RuleID:         rule.ID,
			Enforcement:    rule.Enforcement,
			Scope:          style.Scope(rule.Scope),
			RequirementIDs: append([]string{}, rule.RequirementIDs...),
		})
	}

	return rules
}

func encodeRules(rules []policy.RuleBinding) (schemas []schemaRuleBinding) {
	schemas = make([]schemaRuleBinding, 0, len(rules))
	for _, rule := range rules {
		schemas = append(schemas, schemaRuleBinding{
			ID:             rule.RuleID,
			Enforcement:    rule.Enforcement,
			Scope:          string(rule.Scope),
			RequirementIDs: append([]string{}, rule.RequirementIDs...),
		})
	}

	return schemas
}
