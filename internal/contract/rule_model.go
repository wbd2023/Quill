package contract

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	RuleGroupControlPlane RuleGroup = "control_plane"
	RuleGroupLanguage     RuleGroup = "language_backends"
	RuleGroupRepository   RuleGroup = "repository_scanners"
	RuleGroupExternal     RuleGroup = "external_tools"
)

const (
	ExecutorToolchain      = "toolchain"
	ExecutorControlPlane   = "control_plane"
	ExecutorFileCommand    = "file_command"
	ExecutorGoFormat       = "go_format"
	ExecutorGolangci       = "golangci"
	ExecutorGoStyle        = "go_style"
	ExecutorRepositoryScan = "repository_scan"
)

/* -------------------------------------------- Types ------------------------------------------- */

type RuleGroup string

type Definitions struct {
	Tools []Tool
	Rules []RuleDefinition
}

type RuleDefinition struct {
	ID                  string
	Name                string
	Group               RuleGroup
	Spec                ExecutionSpec
	FixSpec             ExecutionSpec
	RequiredConfigRefs  []string
	RequiredPathClasses []string
}

type Rule struct {
	RuleDefinition
	Level          Level
	Scope          Scope
	RequirementIDs []string
	ConfigRef      string
}

type ExecutionSpec struct {
	Executor       string
	ToolIDs        []string
	ToolID         string
	Check          string
	Scanner        string
	Backend        string
	Language       string
	FileSet        string
	Arguments      []string
	ConfigArgument string
	ConfigFile     string
}

/* ------------------------------------------ Tool IDs ------------------------------------------ */

func (rule RuleDefinition) ToolIDs() (toolIDs []string) {
	return rule.Spec.RequiredToolIDs()
}

func (rule RuleDefinition) FixToolIDs() (toolIDs []string) {
	return rule.FixSpec.RequiredToolIDs()
}

func (spec ExecutionSpec) RequiredToolIDs() (toolIDs []string) {
	if spec.Executor == "" {
		return nil
	}

	if len(spec.ToolIDs) > 0 {
		return append([]string{}, spec.ToolIDs...)
	}

	if spec.ToolID == "" {
		return nil
	}

	return []string{spec.ToolID}
}

func (rule Rule) ToolIDs() (toolIDs []string) {
	return rule.RuleDefinition.ToolIDs()
}

func (rule Rule) FixToolIDs() (toolIDs []string) {
	return rule.RuleDefinition.FixToolIDs()
}
