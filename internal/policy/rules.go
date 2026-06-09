package policy

import "ciphera/tools/internal/style"

// RuleBinding binds a rule capability to scope, requirements, and optional profile inputs.
type RuleBinding struct {
	RuleID         string
	Enforcement    style.Enforcement
	Scope          style.Scope
	RequirementIDs []string
}
