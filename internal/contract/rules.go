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
	ID                 string
	Name               string
	Group              RuleGroup
	Spec               ExecutionSpec
	FixSpec            ExecutionSpec
	RequiredConfigRefs []string
}

type Rule struct {
	ID                 string
	Name               string
	Group              RuleGroup
	Spec               ExecutionSpec
	FixSpec            ExecutionSpec
	RequiredConfigRefs []string
	Level              Level
	Scope              Scope
	RequirementIDs     []string
	ConfigRef          string
	PathClasses        []string
}

func (rule RuleDefinition) ToolIDs() (toolIDs []string) {
	return rule.Spec.RequiredToolIDs()
}

func (rule RuleDefinition) FixToolIDs() (toolIDs []string) {
	return rule.FixSpec.RequiredToolIDs()
}

func (rule Rule) ToolIDs() (toolIDs []string) {
	return rule.Spec.RequiredToolIDs()
}

func (rule Rule) FixToolIDs() (toolIDs []string) {
	return rule.FixSpec.RequiredToolIDs()
}
