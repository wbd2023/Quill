package policy

import "ciphera/tools/internal/contract"

// RuleBinding binds a rule capability to scope, requirements, and optional profile inputs.
type RuleBinding struct {
	RuleID         string
	Enforcement    contract.Enforcement
	Scope          contract.Scope
	RequirementIDs []string
}
