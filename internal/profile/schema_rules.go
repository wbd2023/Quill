package profile

import "ciphera/tools/internal/contract"

type schemaRuleBinding struct {
	RuleID         string         `toml:"rule_id"`
	Level          contract.Level `toml:"level"`
	Scope          string         `toml:"scope"`
	RequirementIDs []string       `toml:"requirement_ids"`
	ConfigRef      string         `toml:"config_ref"`
	Backends       []string       `toml:"backends"`
	PathClasses    []string       `toml:"path_classes"`
}
