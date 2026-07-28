package style

import (
	"slices"
	"time"
)

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
// to inspect requirements, then Bind to produce a Job for the execution. Each Template carries the
// PackID of the Pack that declared it so runtime binding identities stay Pack-qualified and results
// carry provenance.
type Template interface {
	isTemplate()
	describe() (requirements Requirements)
	bind(targets []string) (job Job)
	stamp(packID string) (stamped Template)
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

// StampPackID returns a copy of template with its PackID set to packID. The Pack stamps provenance
// onto every rule execution it declares at registry build time.
func StampPackID(template Template, packID string) (stamped Template) {
	return template.stamp(packID)
}

/* --------------------------------------- Execution Types -------------------------------------- */

// ToolchainExecution represents a check that verifies pinned external tools are installed.
type ToolchainExecution struct {
	PackID  string
	ToolIDs []string
}

// ProfileExecution represents a check that validates the profile configuration.
type ProfileExecution struct {
	PackID string
	Check  string
}

// FileCommandExecution represents running a tool against files selected by a file set.
type FileCommandExecution struct {
	PackID  string
	ToolID  string
	FileSet string

	Arguments []string

	ConfigArgument string
	ConfigFile     string
}

// RepositoryScanExecution represents a repository-wide scan over files from a file set.
type RepositoryScanExecution struct {
	PackID  string
	Scanner string
	FileSet string
}

// TargetCommandTemplate represents running a tool against language-specific targets before target
// resolution.
type TargetCommandTemplate struct {
	PackID  string
	ToolIDs []string

	Action   string
	Language string
}

// TargetCheckTemplate represents a language-specific check before target resolution.
type TargetCheckTemplate struct {
	PackID  string
	ToolIDs []string

	Check    string
	Language string
}

// TargetCommandJob represents a tool run against resolved language-specific targets.
type TargetCommandJob struct {
	PackID  string
	ToolIDs []string

	Action   string
	Language string
	Targets  []string
}

// TargetCheckJob represents a language-specific check against resolved targets.
type TargetCheckJob struct {
	PackID  string
	ToolIDs []string

	Check    string
	Language string
	Targets  []string
}

// ExternalCheckTemplate represents an external Pack check executed as a subprocess over a file set.
// The external protocol is check-only for this MVP: external rules carry no fix job. The template
// is self-describing: it carries the runtime executable command, the Pack directory used to
// resolve it, and the timeout, so one flat driver runs every external Pack check without a
// Pack-qualified binding registry. PackID is stamped by the Pack that declared the rule.
type ExternalCheckTemplate struct {
	PackID        string
	RuleID        string
	CheckID       string
	FileSet       string
	PackDirectory string
	Command       string
	Timeout       time.Duration
}

// ExternalCheckJob is a bound external Pack check ready for execution. It carries everything the
// flat external driver needs to build the request and launch the subprocess.
type ExternalCheckJob struct {
	PackID        string
	RuleID        string
	CheckID       string
	FileSet       string
	PackDirectory string
	Command       string
	Timeout       time.Duration
}

/* -------------------------------------- Interface Methods ------------------------------------- */

func (ToolchainExecution) isTemplate()      {}
func (ProfileExecution) isTemplate()        {}
func (FileCommandExecution) isTemplate()    {}
func (RepositoryScanExecution) isTemplate() {}
func (TargetCommandTemplate) isTemplate()   {}
func (TargetCheckTemplate) isTemplate()     {}
func (ExternalCheckTemplate) isTemplate()   {}

func (e ToolchainExecution) stamp(packID string) (stamped Template) {
	e.PackID = packID
	return e
}

func (e ProfileExecution) stamp(packID string) (stamped Template) {
	e.PackID = packID
	return e
}

func (e FileCommandExecution) stamp(packID string) (stamped Template) {
	e.PackID = packID
	return e
}

func (e RepositoryScanExecution) stamp(packID string) (stamped Template) {
	e.PackID = packID
	return e
}

func (e TargetCommandTemplate) stamp(packID string) (stamped Template) {
	e.PackID = packID
	return e
}

func (e TargetCheckTemplate) stamp(packID string) (stamped Template) {
	e.PackID = packID
	return e
}

func (e ExternalCheckTemplate) stamp(packID string) (stamped Template) {
	e.PackID = packID
	return e
}

func (ToolchainExecution) isJob()      {}
func (ProfileExecution) isJob()        {}
func (FileCommandExecution) isJob()    {}
func (RepositoryScanExecution) isJob() {}
func (TargetCommandJob) isJob()        {}
func (TargetCheckJob) isJob()          {}
func (ExternalCheckJob) isJob()        {}

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

func (e ExternalCheckTemplate) describe() (requirements Requirements) {
	return Requirements{FileSet: e.FileSet}
}

func (e ToolchainExecution) bind([]string) (job Job)      { return e }
func (e ProfileExecution) bind([]string) (job Job)        { return e }
func (e FileCommandExecution) bind([]string) (job Job)    { return e }
func (e RepositoryScanExecution) bind([]string) (job Job) { return e }

func (e ExternalCheckTemplate) bind([]string) (job Job) {
	return ExternalCheckJob(e)
}

func (e TargetCommandTemplate) bind(targets []string) (job Job) {
	return TargetCommandJob{
		PackID:   e.PackID,
		ToolIDs:  slices.Clone(e.ToolIDs),
		Action:   e.Action,
		Language: e.Language,
		Targets:  slices.Clone(targets),
	}
}

func (e TargetCheckTemplate) bind(targets []string) (job Job) {
	return TargetCheckJob{
		PackID:   e.PackID,
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

func (ExternalCheckJob) toolIDs() (ids []string) { return nil }
