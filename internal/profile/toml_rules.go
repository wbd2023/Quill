package profile

import "github.com/wbd2023/quill/internal/style"

type schemaRuleBinding struct {
	ID             string            `toml:"id"`
	Enforcement    style.Enforcement `toml:"enforcement"`
	Scope          string            `toml:"scope"`
	RequirementIDs []string          `toml:"requirement_ids"`
}

func decodeRules(schemas []schemaRuleBinding) (rules []RuleBinding) {
	rules = make([]RuleBinding, 0, len(schemas))
	for _, rule := range schemas {
		rules = append(rules, RuleBinding{
			RuleID:         rule.ID,
			Enforcement:    rule.Enforcement,
			Scope:          style.Scope(rule.Scope),
			RequirementIDs: append([]string{}, rule.RequirementIDs...),
		})
	}

	return rules
}

func encodeRules(rules []RuleBinding) (schemas []schemaRuleBinding) {
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
