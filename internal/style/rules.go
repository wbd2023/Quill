package style

// RuleGroup represents a rule category, for example Go syntax or repository text.
type RuleGroup string

// Definitions represents raw Tool and Rule definitions assembled from Packs
// before Profile compilation produces a Plan.
type Definitions struct {
	ToolIDs []string
	Rules   []RuleDefinition
}

// Plan represents a compiled execution plan: concrete rules with bound enforcement levels, scopes,
// and execution jobs.
type Plan struct {
	Rules []Rule
}

// RuleDefinition represents a pack-declared rule before profile binding. It carries check and
// fix Templates but not enforcement or scope. PackID records the Pack that declared the rule so
// compiled rules and results carry Pack provenance; execution values carry no Pack or rule
// identity of their own.
type RuleDefinition struct {
	ID     string
	PackID string
	Name   string
	Group  RuleGroup

	Check Template
	Fix   Template
}

// Rule represents a Profile-bound executable capability with bound enforcement,
// Scope, Requirements, and execution Jobs. PackID carries the declaring Pack
// through to results.
type Rule struct {
	ID     string
	PackID string
	Name   string
	Group  RuleGroup

	Enforcement    Enforcement
	Scope          Scope
	RequirementIDs []string

	Check Job
	Fix   Job
}

// CheckToolIDs returns the tool IDs required by the rule definition's check Template.
func (rule RuleDefinition) CheckToolIDs() (toolIDs []string) {
	if rule.Check == nil {
		return nil
	}
	return Describe(rule.Check).ToolIDs
}

// FixToolIDs returns the tool IDs required by the rule definition's fix Template.
func (rule RuleDefinition) FixToolIDs() (toolIDs []string) {
	if rule.Fix == nil {
		return nil
	}
	return Describe(rule.Fix).ToolIDs
}

// CheckToolIDs returns the tool IDs required by the rule's check Job.
func (rule Rule) CheckToolIDs() (toolIDs []string) {
	if rule.Check == nil {
		return nil
	}
	return ToolIDs(rule.Check)
}

// FixToolIDs returns the tool IDs required by the rule's fix Job.
func (rule Rule) FixToolIDs() (toolIDs []string) {
	if rule.Fix == nil {
		return nil
	}
	return ToolIDs(rule.Fix)
}
