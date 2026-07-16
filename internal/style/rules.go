package style

// RuleGroup represents a rule category, for example Go syntax or repository text.
type RuleGroup string

// Definitions represents raw tool and rule definitions assembled from packs before the effective
// profile is compiled.
type Definitions struct {
	ToolIDs []string
	Rules   []RuleDefinition
}

// Plan represents a compiled execution plan: concrete rules with bound enforcement levels, scopes,
// and execution jobs.
type Plan struct {
	Rules []Rule
}

// RuleDefinition represents a pack-declared rule before profile binding. It carries check and fix
// execution templates but not enforcement or scope.
type RuleDefinition struct {
	ID    string
	Name  string
	Group RuleGroup

	Check Template
	Fix   Template
}

// Rule represents a profile-bound, enforceable style requirement with bound execution jobs.
type Rule struct {
	ID    string
	Name  string
	Group RuleGroup

	Enforcement    Enforcement
	Scope          Scope
	RequirementIDs []string

	Check Job
	Fix   Job
}

// CheckToolIDs returns the tool IDs required by the rule's check template.
func (rule RuleDefinition) CheckToolIDs() (toolIDs []string) {
	if rule.Check == nil {
		return nil
	}
	return Describe(rule.Check).ToolIDs
}

// FixToolIDs returns the tool IDs required by the rule's fix template.
func (rule RuleDefinition) FixToolIDs() (toolIDs []string) {
	if rule.Fix == nil {
		return nil
	}
	return Describe(rule.Fix).ToolIDs
}

// CheckToolIDs returns the tool IDs required by the rule's check job.
func (rule Rule) CheckToolIDs() (toolIDs []string) {
	if rule.Check == nil {
		return nil
	}
	return ToolIDs(rule.Check)
}

// FixToolIDs returns the tool IDs required by the rule's fix job.
func (rule Rule) FixToolIDs() (toolIDs []string) {
	if rule.Fix == nil {
		return nil
	}
	return ToolIDs(rule.Fix)
}
