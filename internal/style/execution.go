package style

import (
	"slices"
	"time"
)

/* -------------------------------------- Core Abstractions ------------------------------------- */

// TemplateRequirements describes what a Template needs to compile into a runnable Job. It is the
// closed query the profile compiler uses to resolve file sets and targets; drivers switch on the
// concrete Job directly rather than reading TemplateRequirements.
type TemplateRequirements struct {
	ToolIDs []string
	FileSet string

	NeedsTargets    bool
	TargetLanguage  string
	NeedsCheckPaths bool
}

// Template is a sealed, Pack-declared check or fix shape before Profile target binding. Each
// closed variant describes one kind of check or fix a Pack may declare; the profile compiler binds
// a target Template into a Job by resolving its targets. Rule and Pack provenance lives on
// RuleDefinition and Rule, not on a Template or Job: an execution value carries only the strategy
// a binding or driver executes.
type Template interface {
	describe() (requirements TemplateRequirements)
}

// Job is a Profile-compiled check or fix ready for driver dispatch. Non-target values satisfy Job
// directly because Profile binding adds nothing to them; target Jobs are produced by their paired
// Template's Bind method, which copies tool IDs and target names so the Job owns them independently
// of the resolver.
type Job interface {
	toolIDs() (ids []string)
}

// Describe returns the compile-time requirements of a Template. It defensively clones owned slices
// so callers cannot mutate the Template through the returned TemplateRequirements.
func Describe(template Template) (requirements TemplateRequirements) {
	return template.describe()
}

// ToolIDs returns the tool IDs a Job requires. It defensively clones the slice so callers cannot
// mutate the Job through the returned value.
func ToolIDs(job Job) (ids []string) {
	return job.toolIDs()
}

/* ------------------------------------- Non-Target Families ------------------------------------ */

// ToolchainCheck represents a check that verifies pinned external tools are installed. It is
// runnable unchanged once a Rule supplies its policy, so it satisfies both Template and Job.
type ToolchainCheck struct {
	ToolIDs []string
}

// ProfileCheck represents a check that validates the profile configuration. It is runnable
// unchanged once a Rule supplies its policy, so it satisfies both Template and Job.
type ProfileCheck struct {
	Check string
}

// FileCommand represents running a tool against files selected by a file set. It is runnable
// unchanged once a Rule supplies its policy, so it satisfies both Template and Job.
type FileCommand struct {
	ToolID  string
	FileSet string

	Arguments []string

	ConfigArgument string
	ConfigFile     string
}

// RepositoryScan represents a repository-wide scan over files from a file set. It is runnable
// unchanged once a Rule supplies its policy, so it satisfies both Template and Job.
type RepositoryScan struct {
	Scanner string
	FileSet string
}

// ExternalCheck represents an external Pack check executed as a subprocess over a file set. The
// external protocol is check-only for this MVP: external rules carry no fix. The value is
// self-describing: it carries the runtime executable command, the Pack directory used to resolve
// it, and the timeout, so one flat driver runs every external Pack check without a Pack-qualified
// binding registry. Pack and rule provenance for the external request comes from the Rule.
type ExternalCheck struct {
	CheckID       string
	FileSet       string
	PackDirectory string
	Command       string
	Timeout       time.Duration
}

/* -------------------------------------- Target Templates -------------------------------------- */

// TargetCommandTemplate represents running a tool against language-specific targets before Profile
// target binding. Bind resolves it into a TargetCommandJob carrying the inferred target names.
type TargetCommandTemplate struct {
	ToolIDs []string

	Action   string
	Language string
}

// TargetCheckTemplate represents a language-specific check before Profile target binding. Bind
// resolves it into a TargetCheckJob carrying the inferred target names.
type TargetCheckTemplate struct {
	ToolIDs []string

	Check    string
	Language string
}

/* ----------------------------------------- Target Jobs ---------------------------------------- */

// TargetCommandJob is a bound target command: the declared shape plus the resolved target names.
type TargetCommandJob struct {
	ToolIDs []string

	Action   string
	Language string
	Targets  []string
}

// TargetCheckJob is a bound target check: the declared shape plus the resolved target names.
type TargetCheckJob struct {
	ToolIDs []string

	Check    string
	Language string
	Targets  []string
}

/* ---------------------------------------- Phase Methods --------------------------------------- */

func (check ToolchainCheck) describe() (requirements TemplateRequirements) {
	return TemplateRequirements{ToolIDs: slices.Clone(check.ToolIDs)}
}

func (check ToolchainCheck) toolIDs() (ids []string) {
	return slices.Clone(check.ToolIDs)
}

func (check ProfileCheck) describe() (requirements TemplateRequirements) {
	return TemplateRequirements{}
}

func (check ProfileCheck) toolIDs() (ids []string) {
	return nil
}

func (command FileCommand) describe() (requirements TemplateRequirements) {
	toolIDs := []string(nil)
	if command.ToolID != "" {
		toolIDs = []string{command.ToolID}
	}
	return TemplateRequirements{ToolIDs: toolIDs, FileSet: command.FileSet}
}

func (command FileCommand) toolIDs() (ids []string) {
	if command.ToolID == "" {
		return nil
	}
	return []string{command.ToolID}
}

func (scan RepositoryScan) describe() (requirements TemplateRequirements) {
	return TemplateRequirements{FileSet: scan.FileSet}
}

func (scan RepositoryScan) toolIDs() (ids []string) {
	return nil
}

func (check ExternalCheck) describe() (requirements TemplateRequirements) {
	return TemplateRequirements{FileSet: check.FileSet}
}

func (check ExternalCheck) toolIDs() (ids []string) {
	return nil
}

func (template TargetCommandTemplate) describe() (requirements TemplateRequirements) {
	return TemplateRequirements{
		ToolIDs:        slices.Clone(template.ToolIDs),
		NeedsTargets:   true,
		TargetLanguage: template.Language,
	}
}

func (template TargetCheckTemplate) describe() (requirements TemplateRequirements) {
	return TemplateRequirements{
		ToolIDs:         slices.Clone(template.ToolIDs),
		NeedsTargets:    true,
		TargetLanguage:  template.Language,
		NeedsCheckPaths: true,
	}
}

func (job TargetCommandJob) toolIDs() (ids []string) {
	return slices.Clone(job.ToolIDs)
}

func (job TargetCheckJob) toolIDs() (ids []string) {
	return slices.Clone(job.ToolIDs)
}

/* --------------------------------------- Target Binding --------------------------------------- */

// Bind resolves a target command Template into a runnable Job by attaching the inferred target
// names. It clones the tool IDs and target names so the returned Job owns them independently of
// the resolver.
func (template TargetCommandTemplate) Bind(targets []string) (job TargetCommandJob) {
	return TargetCommandJob{
		ToolIDs:  slices.Clone(template.ToolIDs),
		Action:   template.Action,
		Language: template.Language,
		Targets:  slices.Clone(targets),
	}
}

// Bind resolves a target check Template into a runnable Job by attaching the inferred target
// names. It clones the tool IDs and target names so the returned Job owns them independently of
// the resolver.
func (template TargetCheckTemplate) Bind(targets []string) (job TargetCheckJob) {
	return TargetCheckJob{
		ToolIDs:  slices.Clone(template.ToolIDs),
		Check:    template.Check,
		Language: template.Language,
		Targets:  slices.Clone(targets),
	}
}
