package policy

import "ciphera/tools/internal/contract"

type RuleBinding struct {
	RuleID         string
	Level          contract.Level
	Scope          contract.Scope
	RequirementIDs []string
	ConfigRef      string
	Backends       []string
	PathClasses    []string
}
