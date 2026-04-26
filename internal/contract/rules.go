package contract

/* -------------------------------------------- Types ------------------------------------------- */

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

type ExecutionSpec struct {
	Kind   ExecutorKind
	Detail ExecutionDetail
}

type ExecutionDetail interface {
	executionDetail()
}

type ToolchainExecution struct {
	ToolIDs []string
}

type ControlPlaneExecution struct {
	Check string
}

type FileCommandExecution struct {
	ToolID         string
	FileSet        string
	Arguments      []string
	ConfigArgument string
	ConfigFile     string
}

type BackendCommandExecution struct {
	ToolIDs  []string
	Action   string
	Language string
	Backends []string
}

type BackendCheckExecution struct {
	ToolIDs  []string
	Check    string
	Language string
	Backends []string
}

type RepositoryScanExecution struct {
	Scanner string
	FileSet string
}

func (ToolchainExecution) executionDetail() {}

func (ControlPlaneExecution) executionDetail() {}

func (FileCommandExecution) executionDetail() {}

func (BackendCommandExecution) executionDetail() {}

func (BackendCheckExecution) executionDetail() {}

func (RepositoryScanExecution) executionDetail() {}

/* ------------------------------------------ Tool IDs ------------------------------------------ */

func (rule RuleDefinition) ToolIDs() (toolIDs []string) {
	return rule.Spec.RequiredToolIDs()
}

func (rule RuleDefinition) FixToolIDs() (toolIDs []string) {
	return rule.FixSpec.RequiredToolIDs()
}

func (spec ExecutionSpec) RequiredToolIDs() (toolIDs []string) {
	switch detail := spec.Detail.(type) {
	case ToolchainExecution:
		return append([]string{}, detail.ToolIDs...)

	case FileCommandExecution:
		if detail.ToolID == "" {
			return nil
		}
		return []string{detail.ToolID}

	case BackendCommandExecution:
		return append([]string{}, detail.ToolIDs...)

	case BackendCheckExecution:
		return append([]string{}, detail.ToolIDs...)

	default:
		return nil
	}
}

func (rule Rule) ToolIDs() (toolIDs []string) {
	return rule.Spec.RequiredToolIDs()
}

func (rule Rule) FixToolIDs() (toolIDs []string) {
	return rule.FixSpec.RequiredToolIDs()
}

func (spec ExecutionSpec) Empty() (empty bool) {
	return spec.Kind == "" && spec.Detail == nil
}

func (spec ExecutionSpec) Executor() (executor string) {
	return string(spec.Kind)
}

func (spec ExecutionSpec) ToolchainExecution() (detail ToolchainExecution, found bool) {
	detail, found = spec.Detail.(ToolchainExecution)
	return detail, found
}

func (spec ExecutionSpec) ControlPlaneExecution() (detail ControlPlaneExecution, found bool) {
	detail, found = spec.Detail.(ControlPlaneExecution)
	return detail, found
}

func (spec ExecutionSpec) FileCommandExecution() (detail FileCommandExecution, found bool) {
	detail, found = spec.Detail.(FileCommandExecution)
	return detail, found
}

func (spec ExecutionSpec) BackendCommandExecution() (detail BackendCommandExecution, found bool) {
	detail, found = spec.Detail.(BackendCommandExecution)
	return detail, found
}

func (spec ExecutionSpec) BackendCheckExecution() (detail BackendCheckExecution, found bool) {
	detail, found = spec.Detail.(BackendCheckExecution)
	return detail, found
}

func (spec ExecutionSpec) RepositoryScanExecution() (detail RepositoryScanExecution, found bool) {
	detail, found = spec.Detail.(RepositoryScanExecution)
	return detail, found
}

func (spec ExecutionSpec) FileSetName() (name string) {
	switch detail := spec.Detail.(type) {
	case FileCommandExecution:
		return detail.FileSet
	case RepositoryScanExecution:
		return detail.FileSet
	default:
		return ""
	}
}

func (spec ExecutionSpec) BackendLanguage() (language string) {
	switch detail := spec.Detail.(type) {
	case BackendCommandExecution:
		return detail.Language
	case BackendCheckExecution:
		return detail.Language
	}

	return ""
}

func (spec ExecutionSpec) Backends() (backends []string) {
	switch detail := spec.Detail.(type) {
	case BackendCommandExecution:
		return append([]string{}, detail.Backends...)
	case BackendCheckExecution:
		return append([]string{}, detail.Backends...)
	}

	return nil
}

func (spec ExecutionSpec) WithBackends(backends []string) (bound ExecutionSpec) {
	bound = spec
	switch detail := spec.Detail.(type) {
	case BackendCommandExecution:
		detail.Backends = append([]string{}, backends...)
		bound.Detail = detail

	case BackendCheckExecution:
		detail.Backends = append([]string{}, backends...)
		bound.Detail = detail
	}

	return bound
}
