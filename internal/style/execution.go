package style

import "slices"

/* -------------------------------------- Core Abstractions ------------------------------------- */

// Requirements describes what a template needs to compile into a job.
type Requirements struct {
	ToolIDs []string
	FileSet string

	NeedsTargets    bool
	TargetLanguage  string
	NeedsCheckPaths bool
}

// Template is an unbound execution strategy declared by a pack. The profile compiler calls Describe
// to inspect requirements, then Bind to produce a Job for the execution.
type Template interface {
	isTemplate()
	describe() (requirements Requirements)
	bind(targets []string) (job Job)
}

// Job is a bound execution ready for the execution.
type Job interface {
	isJob()
	toolIDs() (ids []string)
}

// Describe returns the requirements of a template.
func Describe(template Template) (requirements Requirements) {
	return template.describe()
}

// Bind resolves targets into a bound job.
func Bind(template Template, targets []string) (job Job) {
	return template.bind(targets)
}

// ToolIDs returns the tool IDs a job requires.
func ToolIDs(job Job) (ids []string) {
	return job.toolIDs()
}

/* --------------------------------------- Execution Types -------------------------------------- */

// ToolchainExecution represents a check that verifies pinned external tools are installed.
type ToolchainExecution struct {
	ToolIDs []string
}

// ProfileExecution represents a check that validates the profile configuration.
type ProfileExecution struct {
	Check string
}

// FileCommandExecution represents running a tool against files selected by a file set.
type FileCommandExecution struct {
	ToolID  string
	FileSet string

	Arguments []string

	ConfigArgument string
	ConfigFile     string
}

// RepositoryScanExecution represents a repository-wide scan over files from a file set.
type RepositoryScanExecution struct {
	Scanner string
	FileSet string
}

// TargetCommandTemplate represents running a tool against language-specific targets before target
// resolution.
type TargetCommandTemplate struct {
	ToolIDs []string

	Action   string
	Language string
}

// TargetCheckTemplate represents a language-specific check before target resolution.
type TargetCheckTemplate struct {
	ToolIDs []string

	Check    string
	Language string
}

// TargetCommandJob represents a tool run against resolved language-specific targets.
type TargetCommandJob struct {
	ToolIDs []string

	Action   string
	Language string
	Targets  []string
}

// TargetCheckJob represents a language-specific check against resolved targets.
type TargetCheckJob struct {
	ToolIDs []string

	Check    string
	Language string
	Targets  []string
}

/* -------------------------------------- Interface Methods ------------------------------------- */

func (ToolchainExecution) isTemplate()      {}
func (ProfileExecution) isTemplate()        {}
func (FileCommandExecution) isTemplate()    {}
func (RepositoryScanExecution) isTemplate() {}
func (TargetCommandTemplate) isTemplate()   {}
func (TargetCheckTemplate) isTemplate()     {}

func (ToolchainExecution) isJob()      {}
func (ProfileExecution) isJob()        {}
func (FileCommandExecution) isJob()    {}
func (RepositoryScanExecution) isJob() {}
func (TargetCommandJob) isJob()        {}
func (TargetCheckJob) isJob()          {}

func (e ToolchainExecution) describe() (requirements Requirements) {
	return Requirements{ToolIDs: slices.Clone(e.ToolIDs)}
}

func (ProfileExecution) describe() (requirements Requirements) {
	return Requirements{}
}

func (e FileCommandExecution) describe() (requirements Requirements) {
	toolIDs := []string(nil)
	if e.ToolID != "" {
		toolIDs = []string{e.ToolID}
	}
	return Requirements{ToolIDs: toolIDs, FileSet: e.FileSet}
}

func (e RepositoryScanExecution) describe() (requirements Requirements) {
	return Requirements{FileSet: e.FileSet}
}

func (e TargetCommandTemplate) describe() (requirements Requirements) {
	return Requirements{
		ToolIDs:        slices.Clone(e.ToolIDs),
		NeedsTargets:   true,
		TargetLanguage: e.Language,
	}
}

func (e TargetCheckTemplate) describe() (requirements Requirements) {
	return Requirements{
		ToolIDs:         slices.Clone(e.ToolIDs),
		NeedsTargets:    true,
		TargetLanguage:  e.Language,
		NeedsCheckPaths: true,
	}
}

func (e ToolchainExecution) bind([]string) (job Job)      { return e }
func (e ProfileExecution) bind([]string) (job Job)        { return e }
func (e FileCommandExecution) bind([]string) (job Job)    { return e }
func (e RepositoryScanExecution) bind([]string) (job Job) { return e }

func (e TargetCommandTemplate) bind(targets []string) (job Job) {
	return TargetCommandJob{
		ToolIDs:  slices.Clone(e.ToolIDs),
		Action:   e.Action,
		Language: e.Language,
		Targets:  slices.Clone(targets),
	}
}

func (e TargetCheckTemplate) bind(targets []string) (job Job) {
	return TargetCheckJob{
		ToolIDs:  slices.Clone(e.ToolIDs),
		Check:    e.Check,
		Language: e.Language,
		Targets:  slices.Clone(targets),
	}
}

func (e ToolchainExecution) toolIDs() (ids []string) { return slices.Clone(e.ToolIDs) }

func (ProfileExecution) toolIDs() (ids []string) { return nil }

func (e FileCommandExecution) toolIDs() (ids []string) {
	if e.ToolID == "" {
		return nil
	}
	return []string{e.ToolID}
}

func (e TargetCommandJob) toolIDs() (ids []string) { return slices.Clone(e.ToolIDs) }
func (e TargetCheckJob) toolIDs() (ids []string)   { return slices.Clone(e.ToolIDs) }

func (RepositoryScanExecution) toolIDs() (ids []string) { return nil }
