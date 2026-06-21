package style

// RuleGroup categorises rules by concern, for example Go syntax or. repository text.
type RuleGroup string

// ExecutionKind names one of the supported execution families for a rule.
type ExecutionKind string

// Definitions holds the raw tool and rule definitions assembled from Packs. before the Effective
// Profile is compiled.
type Definitions struct {
	Tools []Tool
	Rules []RuleDefinition
}

// EffectiveConfig is the compiled configuration: concrete rules bound to. enforcement levels and
// scopes, plus the tools they require.
type EffectiveConfig struct {
	Tools []Tool
	Rules []Rule
}

// ToolByID returns the tool with the given ID, if present in the config.
func (effective EffectiveConfig) ToolByID(id string) (tool Tool, found bool) {
	for _, tool := range effective.Tools {
		if tool.ID == id {
			return tool, true
		}
	}

	return Tool{}, false
}

// RuleDefinition is a Pack-declared rule before profile binding. It carries. the check and fix
// execution specs but not enforcement or scope.
type RuleDefinition struct {
	ID    string
	Name  string
	Group RuleGroup
	Check ExecutionSpec
	Fix   ExecutionSpec
}

// Rule is a profile-bound, enforceable style requirement with a concrete. enforcement level, scope,
// and requirement IDs.
type Rule struct {
	ID             string
	Name           string
	Group          RuleGroup
	Enforcement    Enforcement
	Scope          Scope
	RequirementIDs []string
	Check          ExecutionSpec
	Fix            ExecutionSpec
}

// CheckToolIDs returns the tool IDs required by the rule's check execution spec.
func (rule RuleDefinition) CheckToolIDs() (toolIDs []string) {
	return rule.Check.RequiredToolIDs()
}

// FixToolIDs returns the tool IDs required by the rule's fix execution spec.
func (rule RuleDefinition) FixToolIDs() (toolIDs []string) {
	return rule.Fix.RequiredToolIDs()
}

// CheckToolIDs returns the tool IDs required by the rule's check execution spec.
func (rule Rule) CheckToolIDs() (toolIDs []string) {
	return rule.Check.RequiredToolIDs()
}

// FixToolIDs returns the tool IDs required by the rule's fix execution spec.
func (rule Rule) FixToolIDs() (toolIDs []string) {
	return rule.Fix.RequiredToolIDs()
}
