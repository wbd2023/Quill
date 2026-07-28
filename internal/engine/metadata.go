package engine

import (
	"context"
	"slices"
	"sort"

	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/policy"
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

/* ------------------------------------------ Overview ------------------------------------------ */

// MetadataSnapshot is the complete non-presentation metadata view of one prepared repository.
// It is the single read-only query shared by the discoverability commands (list and explain).
//
// Metadata is metadata-only: it shares the prepare pipeline (Profile, STYLE.md document, shipped
// Pack registry, and compiled effective profile) and never constructs a runner context, resolves
// driver bindings, or inspects tools. The full shipped catalogue is assembled independently to
// enumerate available inventory; the prepared profile selects the active subset.
type MetadataSnapshot struct {
	Packs       []PackMetadata
	Rules       []RuleMetadata
	Tools       []ToolMetadata
	Scopes      []ScopeMetadata
	PackConfigs policy.PackConfigs
}

// PackMetadata describes one catalogue Pack and whether the prepared profile activates it.
type PackMetadata struct {
	ID       string
	Name     string
	Active   bool
	External bool
	ToolIDs  []string
	RuleIDs  []string
}

// RuleMetadata describes one catalogue Rule and, when active, its compiled binding.
type RuleMetadata struct {
	ID             string
	PackID         string
	Name           string
	Group          style.RuleGroup
	Active         bool
	Enforcement    style.Enforcement
	Scope          style.Scope
	RequirementIDs []string
	HasFix         bool
	Check          ExecutionDetail
	Fix            ExecutionDetail
}

// ToolMetadata describes one catalogue Tool capability and the Packs that reference it.
type ToolMetadata struct {
	ID            string
	Name          string
	Command       string
	PinnedVersion string
	PackIDs       []string
	External      bool
}

// ScopeMetadata describes one repository scope.
type ScopeMetadata struct {
	Name    style.Scope
	Roots   []string
	Default bool
}

// ExecutionDetail is the structured, presentation-free description of one rule execution. Category
// is one of the Execution constants below (or empty when a rule declares no execution).
type ExecutionDetail struct {
	Category string
	ToolIDs  []string
	FileSet  string
	Language string
	Detail   string
}

// Metadata loads the prepared repository snapshot and assembles its discoverability metadata. It
// performs no tool inspection and launches no process.
func (engine *Engine) Metadata(
	operationContext context.Context,
) (snapshot MetadataSnapshot, operationError error) {
	if err := operationContext.Err(); err != nil {
		return MetadataSnapshot{}, err
	}

	prepared, err := engine.prepare(operationContext)
	if err != nil {
		return MetadataSnapshot{}, err
	}
	if err := operationContext.Err(); err != nil {
		return MetadataSnapshot{}, err
	}

	catalog := shipped.ComposeCatalog(prepared.externalPacks)
	availablePacks := catalog.Packs()

	packIDs := make([]string, len(availablePacks))
	for index, definition := range availablePacks {
		packIDs[index] = definition.ID
	}

	// The full catalogue registry enumerates every available Pack's resolved rules and tool
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

	return buildMetadataSnapshot(prepared.profile, available, externalIDs), nil
}

/* ------------------------------------------ Assembly ------------------------------------------ */

func buildMetadataSnapshot(
	effective profile.EffectiveProfile,
	registry pack.Registry,
	external map[string]bool,
) (snapshot MetadataSnapshot) {
	config := effective.Profile

	enabled := make(map[string]bool, len(config.EnabledPacks))
	for _, packID := range config.EnabledPacks {
		enabled[packID] = true
	}

	activeRules := make(map[string]style.Rule, len(effective.Effective.Rules))
	for _, rule := range effective.Effective.Rules {
		activeRules[rule.ID] = rule
	}

	snapshot.Packs = buildPackMetadata(registry.Packs(), enabled, external)
	snapshot.Rules = buildRuleMetadata(registry.Rules(), activeRules)
	snapshot.Tools = buildToolMetadata(registry.ToolCapabilities(), registry.Packs(), config.Tools)
	snapshot.Scopes = buildScopeMetadata(config.Repository)
	snapshot.PackConfigs = config.PackConfigs

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

		packs = append(packs, PackMetadata{
			ID:       definition.ID,
			Name:     definition.Name,
			Active:   enabled[definition.ID],
			External: external[definition.ID],
			ToolIDs:  sortedClone(definition.ToolIDs),
			RuleIDs:  ruleIDs,
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
			Check:  executionDetail(definition.Check),
			Fix:    executionDetail(definition.Fix),
		}

		if bound, found := active[definition.ID]; found {
			metadata.Active = true
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
	pinned policy.PinnedTools,
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
			External:      false,
		})
	}

	sort.Slice(tools, func(i int, j int) bool { return tools[i].ID < tools[j].ID })
	return tools
}

func buildScopeMetadata(repository policy.RepositoryConfig) (scopes []ScopeMetadata) {
	scopes = make([]ScopeMetadata, 0, len(repository.ScopeRoots))
	for name, roots := range repository.ScopeRoots {
		scopes = append(scopes, ScopeMetadata{
			Name:    name,
			Roots:   sortedClone(roots),
			Default: name == repository.DefaultScope,
		})
	}

	sort.Slice(scopes, func(i int, j int) bool { return scopes[i].Name < scopes[j].Name })
	return scopes
}

// executionDetail classifies one declared execution template into structured metadata. It reads
// only the template's own fields; it never binds targets or runs anything.
func executionDetail(template style.Template) (detail ExecutionDetail) {
	if template == nil {
		return ExecutionDetail{}
	}

	switch execution := template.(type) {
	case style.ToolchainExecution:
		return ExecutionDetail{
			Category: ExecutionToolchain,
			ToolIDs:  sortedClone(execution.ToolIDs),
		}

	case style.ProfileExecution:
		return ExecutionDetail{
			Category: ExecutionProfile,
			Detail:   execution.Check,
		}

	case style.FileCommandExecution:
		toolIDs := []string(nil)
		if execution.ToolID != "" {
			toolIDs = []string{execution.ToolID}
		}
		return ExecutionDetail{
			Category: ExecutionFileCommand,
			ToolIDs:  toolIDs,
			FileSet:  execution.FileSet,
		}

	case style.RepositoryScanExecution:
		return ExecutionDetail{
			Category: ExecutionRepositoryScan,
			FileSet:  execution.FileSet,
			Detail:   execution.Scanner,
		}

	case style.TargetCommandTemplate:
		return ExecutionDetail{
			Category: ExecutionTargetCommand,
			ToolIDs:  sortedClone(execution.ToolIDs),
			Language: execution.Language,
			Detail:   execution.Action,
		}

	case style.TargetCheckTemplate:
		return ExecutionDetail{
			Category: ExecutionTargetCheck,
			ToolIDs:  sortedClone(execution.ToolIDs),
			Language: execution.Language,
			Detail:   execution.Check,
		}

	case style.ExternalCheckTemplate:
		return ExecutionDetail{
			Category: ExecutionExternal,
			FileSet:  execution.FileSet,
			Detail:   execution.CheckID,
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
