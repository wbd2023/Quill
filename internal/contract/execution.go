package contract

/* -------------------------------------------- Types ------------------------------------------- */

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

type CommandResult struct {
	ExitCode  int
	TimedOut  bool
	Truncated bool
}

/* ------------------------------------------- Markers ------------------------------------------ */

func (ToolchainExecution) executionDetail() {}

func (ControlPlaneExecution) executionDetail() {}

func (FileCommandExecution) executionDetail() {}

func (BackendCommandExecution) executionDetail() {}

func (BackendCheckExecution) executionDetail() {}

func (RepositoryScanExecution) executionDetail() {}

func (spec ExecutionSpec) Empty() (empty bool) {
	return spec.Kind == "" && spec.Detail == nil
}

func (spec ExecutionSpec) Executor() (executor string) {
	return string(spec.Kind)
}

func (spec ExecutionSpec) ToolchainExecution() (execution ToolchainExecution, found bool) {
	execution, found = spec.Detail.(ToolchainExecution)
	return execution, found
}

func (spec ExecutionSpec) ControlPlaneExecution() (execution ControlPlaneExecution, found bool) {
	execution, found = spec.Detail.(ControlPlaneExecution)
	return execution, found
}

func (spec ExecutionSpec) FileCommandExecution() (execution FileCommandExecution, found bool) {
	execution, found = spec.Detail.(FileCommandExecution)
	return execution, found
}

func (spec ExecutionSpec) BackendCommandExecution() (
	execution BackendCommandExecution,
	found bool,
) {
	execution, found = spec.Detail.(BackendCommandExecution)
	return execution, found
}

func (spec ExecutionSpec) BackendCheckExecution() (execution BackendCheckExecution, found bool) {
	execution, found = spec.Detail.(BackendCheckExecution)
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

	case BackendCommandExecution:
		return append([]string{}, execution.ToolIDs...)

	case BackendCheckExecution:
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

func (spec ExecutionSpec) BackendLanguage() (language string) {
	switch execution := spec.Detail.(type) {
	case BackendCommandExecution:
		return execution.Language
	case BackendCheckExecution:
		return execution.Language
	}

	return ""
}

func (spec ExecutionSpec) Backends() (backends []string) {
	switch execution := spec.Detail.(type) {
	case BackendCommandExecution:
		return append([]string{}, execution.Backends...)
	case BackendCheckExecution:
		return append([]string{}, execution.Backends...)
	}

	return nil
}

func (spec ExecutionSpec) WithBackends(backends []string) (bound ExecutionSpec) {
	bound = spec
	switch execution := spec.Detail.(type) {
	case BackendCommandExecution:
		execution.Backends = append([]string{}, backends...)
		bound.Detail = execution

	case BackendCheckExecution:
		execution.Backends = append([]string{}, backends...)
		bound.Detail = execution
	}

	return bound
}
