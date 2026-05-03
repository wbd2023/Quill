package policy

import "ciphera/tools/internal/contract"

// RulePackConfig selects reusable rule capability packs.
type RulePackConfig struct {
	Enabled []string
}

// RuleBinding binds a rule capability to scope, requirements, and optional profile inputs.
type RuleBinding struct {
	RuleID          string
	Level           contract.Level
	Scope           contract.Scope
	RequirementIDs  []string
	ConfigReference string
	Backends        []string
	PathClasses     []string
}
