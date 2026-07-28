package pack

import (
	"fmt"
	"sort"

	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

/* ------------------------------------------ Registry ------------------------------------------ */

// Registry stores selected Pack definitions as runtime rule and tool definitions. Rules carry the
// PackID of their declaring Pack (provenance), and capabilities are resolved copies of the
// catalogue's canonical Tools.
type Registry struct {
	packs        []Definition
	capabilities []toolchain.Capability
	rules        []style.RuleDefinition
}

// Packs returns defensive copies of the packs registered in the registry.
func (registry Registry) Packs() (packs []Definition) {
	return CloneDefinitions(registry.packs)
}

// ToolCapabilities returns defensive deep copies of the tool capabilities registered in the
// registry, including mutable installer data.
func (registry Registry) ToolCapabilities() (capabilities []toolchain.Capability) {
	return CloneCapabilities(registry.capabilities)
}

// Rules returns defensive copies of the rule definitions registered in the registry.
func (registry Registry) Rules() (rules []style.RuleDefinition) {
	return CloneRules(registry.rules)
}

// Definitions returns the registered tool IDs and rule definitions.
func (registry Registry) Definitions() (definitions style.Definitions) {
	toolIDs := make([]string, len(registry.capabilities))
	for i, capability := range registry.capabilities {
		toolIDs[i] = capability.ID
	}

	return style.Definitions{
		ToolIDs: toolIDs,
		Rules:   registry.Rules(),
	}
}

/* ------------------------------------------ Assembly ------------------------------------------ */

func selectPacks(available []Definition, enabled []string) (selected []Definition, err error) {
	packByID := make(map[string]Definition, len(available))
	for _, pack := range available {
		packByID[pack.ID] = pack
	}

	selected = make([]Definition, 0, len(enabled))
	seen := make(map[string]bool, len(enabled))
	for _, packID := range enabled {
		if seen[packID] {
			return nil, fmt.Errorf("duplicate pack %q", packID)
		}

		pack, found := packByID[packID]
		if !found {
			return nil, fmt.Errorf("unknown pack %q", packID)
		}

		seen[packID] = true
		selected = append(selected, pack)
	}

	return selected, nil
}

// buildRegistry resolves each Pack's Tool references against the canonical capabilities and stamps
// every rule with its declaring Pack's ID. Unknown Tool references are rejected before any rule or
// capability is recorded.
func buildRegistry(
	tools []toolchain.Capability,
	packs []Definition,
) (registry Registry, err error) {
	registry.packs = CloneDefinitions(packs)

	toolByID := make(map[string]toolchain.Capability, len(tools))
	for _, tool := range tools {
		toolByID[tool.ID] = tool
	}

	seenTools := make(map[string]bool)
	orderedToolIDs := make([]string, 0)
	for _, pack := range packs {
		if err = validatePackRuleToolScope(pack); err != nil {
			return Registry{}, err
		}

		for _, toolID := range pack.ToolIDs {
			_, known := toolByID[toolID]
			if !known {
				return Registry{}, fmt.Errorf(
					"pack %q references unknown tool %q",
					pack.ID,
					toolID,
				)
			}

			if seenTools[toolID] {
				continue
			}

			seenTools[toolID] = true
			orderedToolIDs = append(orderedToolIDs, toolID)
		}

		registry.rules = append(registry.rules, stampPackRules(pack)...)
	}

	sort.Strings(orderedToolIDs)
	registry.capabilities = make([]toolchain.Capability, 0, len(orderedToolIDs))
	for _, toolID := range orderedToolIDs {
		registry.capabilities = append(registry.capabilities, toolByID[toolID])
	}

	return registry, nil
}

// stampPackRules returns copies of pack's rules with the Pack's ID carried on each rule and its
// check/fix templates so compiled rules, jobs, and results retain Pack provenance.
func stampPackRules(pack Definition) (rules []style.RuleDefinition) {
	rules = make([]style.RuleDefinition, len(pack.Rules))
	for index, rule := range pack.Rules {
		stamped := rule
		stamped.PackID = pack.ID
		stamped.Check = stampTemplate(rule.Check, pack.ID)
		stamped.Fix = stampTemplate(rule.Fix, pack.ID)
		rules[index] = stamped
	}

	return rules
}

func stampTemplate(template style.Template, packID string) (stamped style.Template) {
	if template == nil {
		return nil
	}

	return style.StampPackID(template, packID)
}

// validatePackRuleToolScope ensures every rule a Pack declares only references Tools the Pack
// itself declares. A rule may not rely on a Tool brought in by another Pack: Tool IDs are global,
// but each Pack owns the Tools its own rules exercise, so a cross-Pack Tool reference is rejected
// at assembly even when the Tool is present in the aggregate catalogue.
func validatePackRuleToolScope(pack Definition) (err error) {
	declared := make(map[string]bool, len(pack.ToolIDs))
	for _, toolID := range pack.ToolIDs {
		declared[toolID] = true
	}

	for _, rule := range pack.Rules {
		for _, label := range []string{"check", "fix"} {
			toolIDs := ruleToolIDs(rule, label)
			for _, toolID := range toolIDs {
				if !declared[toolID] {
					return fmt.Errorf(
						"rule %q in pack %q %s references tool %q not declared by the pack",
						rule.ID,
						pack.ID,
						label,
						toolID,
					)
				}
			}
		}
	}

	return nil
}

func ruleToolIDs(rule style.RuleDefinition, label string) (toolIDs []string) {
	if label == "check" {
		return rule.CheckToolIDs()
	}
	return rule.FixToolIDs()
}

/* ----------------------------------------- Validation ----------------------------------------- */

func validateRegistry(registry Registry) (err error) {
	if err = validatePackFileSets(registry.packs); err != nil {
		return err
	}

	seenRuleIDs := make(map[string]bool, len(registry.rules))
	for _, rule := range registry.rules {
		if rule.ID == "" {
			return fmt.Errorf("pack registry contains an empty rule id")
		}

		if rule.Check == nil {
			return fmt.Errorf("rule %q has no check execution", rule.ID)
		}

		if seenRuleIDs[rule.ID] {
			return fmt.Errorf("duplicate rule id %q", rule.ID)
		}

		seenRuleIDs[rule.ID] = true
	}

	return nil
}

func validatePackFileSets(packs []Definition) (err error) {
	packByFileSet := make(map[string]string)
	for _, pack := range packs {
		for _, fileSet := range pack.FileSets {
			if fileSet.Name == "" {
				return fmt.Errorf("pack %q contains a file set with an empty name", pack.ID)
			}

			owner, found := packByFileSet[fileSet.Name]
			if found {
				return fmt.Errorf(
					"file set %q is defined by both packs %q and %q",
					fileSet.Name,
					owner,
					pack.ID,
				)
			}

			packByFileSet[fileSet.Name] = pack.ID
		}
	}

	return nil
}
