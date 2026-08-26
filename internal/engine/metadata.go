package engine

import (
	"context"
	"slices"
	"sort"

	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

// Execution categories classify a rule's check or fix execution.
const (
	ExecutionToolchain      = "toolchain"
	ExecutionProfile        = "profile"
	ExecutionFileCommand    = "file_command"
	ExecutionRepositoryScan = "repository_scan"
	ExecutionTargetCommand  = "target_command"
	ExecutionTargetCheck    = "target_check"
	ExecutionExternal       = "external"
)

const (
	PackProvenanceShipped  PackProvenance = "shipped"
	PackProvenanceExternal PackProvenance = "external"
)

/* ------------------------------------------ Overview ------------------------------------------ */

// PackProvenance identifies whether a Pack ships with Quill or is repository-provided.
type PackProvenance string

// MetadataSnapshot is the complete non-presentation metadata view of one prepared repository.
// It is the single read-only query shared by the discoverability commands (list and explain).
//
// Metadata is metadata-only: it shares the prepare pipeline (Profile, STYLE.md
// document, Pack catalog, and compiled Plan) and never constructs a runner
// context, resolves driver bindings, or inspects tools. The full catalog is
// assembled independently to enumerate available inventory; the prepared
// Profile selects the active subset.
type MetadataSnapshot struct {
	Packs        []PackMetadata
	Rules        []RuleMetadata
	Tools        []ToolMetadata
	Scopes       []ScopeMetadata
	PackPolicies profile.PackPolicies
}

// PackMetadata describes one catalog Pack and whether the prepared profile activates it.
type PackMetadata struct {
	ID         string
	Name       string
	Enabled    bool
	Provenance PackProvenance
	ToolIDs    []string
	RuleIDs    []string
}

// RuleMetadata describes one catalog Rule and, when active, its compiled binding.
type RuleMetadata struct {
	ID             string
	PackID         string
	Name           string
	Group          style.RuleGroup
	Enabled        bool
	Enforcement    style.Enforcement
	Scope          style.Scope
	RequirementIDs []string
	HasFix         bool
	Check          ExecutionDetail
	Fix            ExecutionDetail
}

// ToolMetadata describes one catalog Tool capability and the Packs that reference it.
type ToolMetadata struct {
	ID            string
	Name          string
	Command       string
	PinnedVersion string
	PackIDs       []string
}

// ScopeMetadata describes one repository scope.
type ScopeMetadata struct {
	Name      style.Scope
	Roots     []string
	IsDefault bool
}

// ExecutionDetail is the structured, presentation-free description of one rule execution. Category
// is one of the Execution constants below (or empty when a rule declares no execution).
type ExecutionDetail struct {
	Category string
	ToolIDs  []string
	FileSet  string
	Language string
}

// Metadata loads the prepared repository snapshot and assembles its discoverability metadata. It
// performs no tool inspection and launches no process.
func (engine *Engine) Metadata(
	ctx context.Context,
) (snapshot MetadataSnapshot, err error) {
	if err := ctx.Err(); err != nil {
		return MetadataSnapshot{}, err
	}

	prepared, err := engine.prepare(ctx)
	if err != nil {
		return MetadataSnapshot{}, err
	}

	if err := ctx.Err(); err != nil {
		return MetadataSnapshot{}, err
	}

	catalog := shipped.ComposeCatalog(prepared.externalPacks)
	availablePacks := catalog.Packs()

	packIDs := make([]string, len(availablePacks))
	for index, definition := range availablePacks {
		packIDs[index] = definition.ID
	}

	// The full catalog registry enumerates every available Pack's resolved rules and tool
	// capabilities with Pack provenance stamped on each rule. It is pure declaration assembly:
	// no driver binding and no tool inspection.
	available, err := catalog.Registry(packIDs)
	if err != nil {
		return MetadataSnapshot{}, err
	}

	externalIDs := make(map[string]bool, len(prepared.externalPacks))
	for _, definition := range prepared.externalPacks {
		externalIDs[definition.ID] = true
	}

	return buildMetadataSnapshot(prepared.config, prepared.plan, available, externalIDs), nil
}

/* ------------------------------------------ Assembly ------------------------------------------ */

func buildMetadataSnapshot(
	config profile.Profile,
	plan style.Plan,
	registry pack.Registry,
	external map[string]bool,
) (snapshot MetadataSnapshot) {
	enabled := make(map[string]bool, len(config.EnabledPacks))
	for _, packID := range config.EnabledPacks {
		enabled[packID] = true
	}

	activeRules := make(map[string]style.Rule, len(plan.Rules))
	for _, rule := range plan.Rules {
		activeRules[rule.ID] = rule
	}

	snapshot.Packs = buildPackMetadata(registry.Packs(), enabled, external)
	snapshot.Rules = buildRuleMetadata(registry.Rules(), activeRules)
	snapshot.Tools = buildToolMetadata(registry.ToolCapabilities(), registry.Packs(), config.Tools)
	snapshot.Scopes = buildScopeMetadata(config.Repository)
	snapshot.PackPolicies = config.PackPolicies

	return snapshot
}

func buildPackMetadata(
	definitions []pack.Definition,
	enabled map[string]bool,
	external map[string]bool,
) (packs []PackMetadata) {
	packs = make([]PackMetadata, 0, len(definitions))
	for _, definition := range definitions {
		ruleIDs := make([]string, len(definition.Rules))
		for index, rule := range definition.Rules {
			ruleIDs[index] = rule.ID
		}
		slices.Sort(ruleIDs)

		provenance := PackProvenanceShipped
		if external[definition.ID] {
			provenance = PackProvenanceExternal
		}

		packs = append(packs, PackMetadata{
			ID:         definition.ID,
			Name:       definition.Name,
			Enabled:    enabled[definition.ID],
			Provenance: provenance,
			ToolIDs:    sortedClone(definition.ToolIDs),
			RuleIDs:    ruleIDs,
		})
	}

	sort.Slice(packs, func(i int, j int) bool { return packs[i].ID < packs[j].ID })
	return packs
}

func buildRuleMetadata(
	definitions []style.RuleDefinition,
	active map[string]style.Rule,
) (rules []RuleMetadata) {
	rules = make([]RuleMetadata, 0, len(definitions))
	for _, definition := range definitions {
		metadata := RuleMetadata{
			ID:     definition.ID,
			PackID: definition.PackID,
			Name:   definition.Name,
			Group:  definition.Group,
			HasFix: definition.Fix != nil,
			Check:  classifyExecutionDetail(definition.Check),
			Fix:    classifyExecutionDetail(definition.Fix),
		}

		if bound, found := active[definition.ID]; found {
			metadata.Enabled = true
			metadata.Enforcement = bound.Enforcement
			metadata.Scope = bound.Scope
			metadata.RequirementIDs = sortedClone(bound.RequirementIDs)
		}

		rules = append(rules, metadata)
	}

	sort.Slice(rules, func(i int, j int) bool { return rules[i].ID < rules[j].ID })
	return rules
}

func buildToolMetadata(
	capabilities []toolchain.Capability,
	packs []pack.Definition,
	pinned profile.PinnedTools,
) (tools []ToolMetadata) {
	ownerPacks := make(map[string][]string)
	for _, definition := range packs {
		for _, toolID := range definition.ToolIDs {
			ownerPacks[toolID] = append(ownerPacks[toolID], definition.ID)
		}
	}

	tools = make([]ToolMetadata, 0, len(capabilities))
	for _, capability := range capabilities {
		packIDs := sortedClone(ownerPacks[capability.ID])

		version := ""
		if pin, found := pinned.Lookup(capability.ID); found {
			version = pin.Version
		}

		tools = append(tools, ToolMetadata{
			ID:            capability.ID,
			Name:          capability.Name,
			Command:       capability.Command,
			PinnedVersion: version,
			PackIDs:       packIDs,
		})
	}

	sort.Slice(tools, func(i int, j int) bool { return tools[i].ID < tools[j].ID })
	return tools
}

func buildScopeMetadata(repository profile.RepositoryConfig) (scopes []ScopeMetadata) {
	scopes = make([]ScopeMetadata, 0, len(repository.ScopeRoots))
	for name, roots := range repository.ScopeRoots {
		scopes = append(scopes, ScopeMetadata{
			Name:      name,
			Roots:     sortedClone(roots),
			IsDefault: name == repository.DefaultScope,
		})
	}

	sort.Slice(scopes, func(i int, j int) bool { return scopes[i].Name < scopes[j].Name })
	return scopes
}

// classifyExecutionDetail classifies one declared Template into structured metadata. It reads
// only the Template's own fields; it never binds targets or runs anything.
func classifyExecutionDetail(template style.Template) (detail ExecutionDetail) {
	if template == nil {
		return ExecutionDetail{}
	}

	switch execution := template.(type) {
	case style.ToolchainCheck:
		return ExecutionDetail{
			Category: ExecutionToolchain,
			ToolIDs:  sortedClone(execution.ToolIDs),
		}

	case style.ProfileCheck:
		return ExecutionDetail{Category: ExecutionProfile}

	case style.FileCommand:
		toolIDs := []string(nil)
		if execution.ToolID != "" {
			toolIDs = []string{execution.ToolID}
		}
		return ExecutionDetail{
			Category: ExecutionFileCommand,
			ToolIDs:  toolIDs,
			FileSet:  execution.FileSet,
		}

	case style.RepositoryScan:
		return ExecutionDetail{
			Category: ExecutionRepositoryScan,
			FileSet:  execution.FileSet,
		}

	case style.TargetCommandTemplate:
		return ExecutionDetail{
			Category: ExecutionTargetCommand,
			ToolIDs:  sortedClone(execution.ToolIDs),
			Language: execution.Language,
		}

	case style.TargetCheckTemplate:
		return ExecutionDetail{
			Category: ExecutionTargetCheck,
			ToolIDs:  sortedClone(execution.ToolIDs),
			Language: execution.Language,
		}

	case style.ExternalCheck:
		return ExecutionDetail{
			Category: ExecutionExternal,
			FileSet:  execution.FileSet,
		}

	default:
		return ExecutionDetail{}
	}
}

func sortedClone(values []string) (clone []string) {
	if len(values) == 0 {
		return nil
	}

	clone = append([]string(nil), values...)
	slices.Sort(clone)
	return clone
}
