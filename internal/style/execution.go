package style

// ExecutionSpec describes how a rule is executed. The concrete Detail type determines which driver
// handles the rule.
type ExecutionSpec struct {
	Detail ExecutionDetail
}

// ExecutionDetail is a sealed interface implemented by each execution detail type. The marker
// method is unexported so only types in this package can satisfy it.
type ExecutionDetail interface {
	executionDetail()
}

// ToolchainExecution validates that pinned external tools are installed and healthy.
type ToolchainExecution struct {
	ToolIDs []string
}

// ProfileExecution validates the profile configuration identified by its check ID.
type ProfileExecution struct {
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

func (ProfileExecution) executionDetail() {}

func (FileCommandExecution) executionDetail() {}

func (TargetCommandExecution) executionDetail() {}

func (TargetCheckExecution) executionDetail() {}

func (RepositoryScanExecution) executionDetail() {}

/* ------------------------------------------ Accessors ----------------------------------------- */

// Empty reports whether the spec has no detail.
func (spec ExecutionSpec) Empty() (empty bool) {
	return spec.Detail == nil
}

// ToolchainExecution returns the toolchain execution detail, if the spec holds one.
func (spec ExecutionSpec) ToolchainExecution() (execution ToolchainExecution, found bool) {
	execution, found = spec.Detail.(ToolchainExecution)
	return execution, found
}

// ProfileExecution returns the profile execution detail, if the spec holds one.
func (spec ExecutionSpec) ProfileExecution() (execution ProfileExecution, found bool) {
	execution, found = spec.Detail.(ProfileExecution)
	return execution, found
}

// FileCommandExecution returns the file-command execution detail, if the spec holds one.
func (spec ExecutionSpec) FileCommandExecution() (execution FileCommandExecution, found bool) {
	execution, found = spec.Detail.(FileCommandExecution)
	return execution, found
}

// TargetCommandExecution returns the target-command execution detail, if the spec holds one.
func (spec ExecutionSpec) TargetCommandExecution() (execution TargetCommandExecution, found bool) {
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
