package style

/* --------------------------------------- Execution Kinds -------------------------------------- */

const (
	ExecutionToolchain      ExecutionKind = "toolchain"
	ExecutionProject        ExecutionKind = "project"
	ExecutionFileCommand    ExecutionKind = "file_command"
	ExecutionTargetCommand  ExecutionKind = "target_command"
	ExecutionTargetCheck    ExecutionKind = "target_check"
	ExecutionRepositoryScan ExecutionKind = "repository_scan"
)

/* -------------------------------------------- Types ------------------------------------------- */

type ExecutionSpec struct {
	Kind   ExecutionKind
	Detail ExecutionDetail
}

type ExecutionDetail interface {
	executionDetail()
}

type ToolchainExecution struct {
	ToolIDs []string
}

type ProjectExecution struct {
	Check string
}

type FileCommandExecution struct {
	ToolID         string
	FileSet        string
	Arguments      []string
	ConfigArgument string
	ConfigFile     string
}

type TargetCommandExecution struct {
	ToolIDs  []string
	Action   string
	Language string
	Targets  []string
}

type TargetCheckExecution struct {
	ToolIDs  []string
	Check    string
	Language string
	Targets  []string
}

type RepositoryScanExecution struct {
	Scanner string
	FileSet string
}

type CommandResult struct {
	ExitCode  int
	TimedOut  bool
	Truncated bool
}

/* ------------------------------------------- Markers ------------------------------------------ */

func (ToolchainExecution) executionDetail() {}

func (ProjectExecution) executionDetail() {}

func (FileCommandExecution) executionDetail() {}

func (TargetCommandExecution) executionDetail() {}

func (TargetCheckExecution) executionDetail() {}

func (RepositoryScanExecution) executionDetail() {}

func (spec ExecutionSpec) Empty() (empty bool) {
	return spec.Kind == "" && spec.Detail == nil
}

func (spec ExecutionSpec) ToolchainExecution() (execution ToolchainExecution, found bool) {
	execution, found = spec.Detail.(ToolchainExecution)
	return execution, found
}

func (spec ExecutionSpec) ProjectExecution() (execution ProjectExecution, found bool) {
	execution, found = spec.Detail.(ProjectExecution)
	return execution, found
}

func (spec ExecutionSpec) FileCommandExecution() (execution FileCommandExecution, found bool) {
	execution, found = spec.Detail.(FileCommandExecution)
	return execution, found
}

func (spec ExecutionSpec) TargetCommandExecution() (
	execution TargetCommandExecution,
	found bool,
) {
	execution, found = spec.Detail.(TargetCommandExecution)
	return execution, found
}

func (spec ExecutionSpec) TargetCheckExecution() (execution TargetCheckExecution, found bool) {
	execution, found = spec.Detail.(TargetCheckExecution)
	return execution, found
}

func (spec ExecutionSpec) RepositoryScanExecution() (
	execution RepositoryScanExecution,
	found bool,
) {
	execution, found = spec.Detail.(RepositoryScanExecution)
	return execution, found
}

/* ------------------------------------------- Queries ------------------------------------------ */

func (spec ExecutionSpec) RequiredToolIDs() (toolIDs []string) {
	switch execution := spec.Detail.(type) {
	case ToolchainExecution:
		return append([]string{}, execution.ToolIDs...)

	case FileCommandExecution:
		if execution.ToolID == "" {
			return nil
		}
		return []string{execution.ToolID}

	case TargetCommandExecution:
		return append([]string{}, execution.ToolIDs...)

	case TargetCheckExecution:
		return append([]string{}, execution.ToolIDs...)

	default:
		return nil
	}
}

func (spec ExecutionSpec) FileSetName() (name string) {
	switch execution := spec.Detail.(type) {
	case FileCommandExecution:
		return execution.FileSet
	case RepositoryScanExecution:
		return execution.FileSet
	default:
		return ""
	}
}

func (spec ExecutionSpec) UsesTargets() (uses bool) {
	switch spec.Detail.(type) {
	case TargetCommandExecution, TargetCheckExecution:
		return true
	default:
		return false
	}
}

func (spec ExecutionSpec) TargetLanguage() (language string) {
	switch execution := spec.Detail.(type) {
	case TargetCommandExecution:
		return execution.Language
	case TargetCheckExecution:
		return execution.Language
	}

	return ""
}

func (spec ExecutionSpec) RequiresTargetCheckPaths() (requires bool) {
	_, requires = spec.Detail.(TargetCheckExecution)
	return requires
}

func (spec ExecutionSpec) Targets() (targets []string) {
	switch execution := spec.Detail.(type) {
	case TargetCommandExecution:
		return append([]string{}, execution.Targets...)
	case TargetCheckExecution:
		return append([]string{}, execution.Targets...)
	}

	return nil
}

func (spec ExecutionSpec) WithTargets(targets []string) (bound ExecutionSpec) {
	bound = spec
	switch execution := spec.Detail.(type) {
	case TargetCommandExecution:
		execution.Targets = append([]string{}, targets...)
		bound.Detail = execution

	case TargetCheckExecution:
		execution.Targets = append([]string{}, targets...)
		bound.Detail = execution
	}

	return bound
}
