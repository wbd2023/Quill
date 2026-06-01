package contract

type RuleGroup string

type ExecutorKind string

type Definitions struct {
	Tools []Tool
	Rules []RuleDefinition
}

type EffectiveConfig struct {
	Tools []Tool
	Rules []Rule
}

func (effective EffectiveConfig) ToolByID(id string) (tool Tool, found bool) {
	for _, tool := range effective.Tools {
		if tool.ID == id {
			return tool, true
		}
	}

	return Tool{}, false
}

type RuleDefinition struct {
	ID    string
	Name  string
	Group RuleGroup
	Check ExecutionSpec
	Fix   ExecutionSpec
}

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

func (rule RuleDefinition) CheckToolIDs() (toolIDs []string) {
	return rule.Check.RequiredToolIDs()
}

func (rule RuleDefinition) FixToolIDs() (toolIDs []string) {
	return rule.Fix.RequiredToolIDs()
}

func (rule Rule) CheckToolIDs() (toolIDs []string) {
	return rule.Check.RequiredToolIDs()
}

func (rule Rule) FixToolIDs() (toolIDs []string) {
	return rule.Fix.RequiredToolIDs()
}
