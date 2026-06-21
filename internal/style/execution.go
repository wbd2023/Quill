package style

/* --------------------------------------- Execution Kinds -------------------------------------- */

// ExecutionKind constants name the supported execution families for a rule.
const (
	// ExecutionToolchain runs a pinned tool against the repository root.
	ExecutionToolchain ExecutionKind = "toolchain"
	// ExecutionProject runs a repository-wide check that does not invoke an external tool.
	ExecutionProject ExecutionKind = "project"
	// ExecutionFileCommand runs a tool against files selected by a file set.
	ExecutionFileCommand ExecutionKind = "file_command"
	// ExecutionTargetCommand runs a tool against language-specific targets.
	ExecutionTargetCommand ExecutionKind = "target_command"
	// ExecutionTargetCheck runs a language-specific check against targets.
	ExecutionTargetCheck ExecutionKind = "target_check"
	// ExecutionRepositoryScan runs a repository-wide scanner that does not invoke an external tool.
	ExecutionRepositoryScan ExecutionKind = "repository_scan"
)

/* -------------------------------------------- Types ------------------------------------------- */

// ExecutionSpec describes how a rule is executed: the execution kind and its
// concrete detail.
type ExecutionSpec struct {
	Kind   ExecutionKind
	Detail ExecutionDetail
}

// ExecutionDetail is a sealed interface implemented by each execution-family. detail type. The
// marker method is unexported so only types in this package can satisfy it.
type ExecutionDetail interface {
	executionDetail()
}

// ToolchainExecution runs one or more pinned tools against the repository root.
type ToolchainExecution struct {
	ToolIDs []string
}

// ProjectExecution runs a repository-wide check identified by its check ID.
type ProjectExecution struct {
	Check string
}

// FileCommandExecution runs a tool against files selected by a file set.
type FileCommandExecution struct {
	ToolID         string
	FileSet        string
	Arguments      []string
	ConfigArgument string
	ConfigFile     string
}

// TargetCommandExecution runs a tool against language-specific targets.
type TargetCommandExecution struct {
	ToolIDs  []string
	Action   string
	Language string
	Targets  []string
}

// TargetCheckExecution runs a language-specific check against targets.
type TargetCheckExecution struct {
	ToolIDs  []string
	Check    string
	Language string
	Targets  []string
}

// RepositoryScanExecution runs a repository-wide scanner over files from a file set.
type RepositoryScanExecution struct {
	Scanner string
	FileSet string
}

// CommandResult holds the raw outcome of an external command execution.
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

// Empty reports whether the spec has no kind and no detail.
func (spec ExecutionSpec) Empty() (empty bool) {
	return spec.Kind == "" && spec.Detail == nil
}

// ToolchainExecution returns the toolchain execution detail, if the spec holds one.
func (spec ExecutionSpec) ToolchainExecution() (execution ToolchainExecution, found bool) {
	execution, found = spec.Detail.(ToolchainExecution)
	return execution, found
}

// ProjectExecution returns the project execution detail, if the spec holds one.
func (spec ExecutionSpec) ProjectExecution() (execution ProjectExecution, found bool) {
	execution, found = spec.Detail.(ProjectExecution)
	return execution, found
}

// FileCommandExecution returns the file-command execution detail, if the spec holds one.
func (spec ExecutionSpec) FileCommandExecution() (execution FileCommandExecution, found bool) {
	execution, found = spec.Detail.(FileCommandExecution)
	return execution, found
}

// TargetCommandExecution returns the target-command execution detail, if the spec holds one.
func (spec ExecutionSpec) TargetCommandExecution() (
	execution TargetCommandExecution,
	found bool,
) {
	execution, found = spec.Detail.(TargetCommandExecution)
	return execution, found
}

// TargetCheckExecution returns the target-check execution detail, if the spec holds one.
func (spec ExecutionSpec) TargetCheckExecution() (execution TargetCheckExecution, found bool) {
	execution, found = spec.Detail.(TargetCheckExecution)
	return execution, found
}

// RepositoryScanExecution returns the repository-scan execution detail, if the spec holds one.
func (spec ExecutionSpec) RepositoryScanExecution() (
	execution RepositoryScanExecution,
	found bool,
) {
	execution, found = spec.Detail.(RepositoryScanExecution)
	return execution, found
}

/* ------------------------------------------- Queries ------------------------------------------ */

// RequiredToolIDs returns the tool IDs the spec needs to execute, or nil if none.
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

// FileSetName returns the file set name used by the spec, or empty if none.
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

// UsesTargets reports whether the spec executes against language-specific targets.
func (spec ExecutionSpec) UsesTargets() (uses bool) {
	switch spec.Detail.(type) {
	case TargetCommandExecution, TargetCheckExecution:
		return true
	default:
		return false
	}
}

// TargetLanguage returns the language for target execution, or empty if none.
func (spec ExecutionSpec) TargetLanguage() (language string) {
	switch execution := spec.Detail.(type) {
	case TargetCommandExecution:
		return execution.Language
	case TargetCheckExecution:
		return execution.Language
	}

	return ""
}

// RequiresTargetCheckPaths reports whether the spec is a target check needing path arguments.
func (spec ExecutionSpec) RequiresTargetCheckPaths() (requires bool) {
	_, requires = spec.Detail.(TargetCheckExecution)
	return requires
}

// Targets returns the target paths bound to the spec, or nil if none.
func (spec ExecutionSpec) Targets() (targets []string) {
	switch execution := spec.Detail.(type) {
	case TargetCommandExecution:
		return append([]string{}, execution.Targets...)
	case TargetCheckExecution:
		return append([]string{}, execution.Targets...)
	}

	return nil
}

// WithTargets returns a copy of the spec with the given targets bound to it.
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
