package pack

import (
	"fmt"
	"sort"

	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

/* ------------------------------------------ Registry ------------------------------------------ */

// Registry stores selected Pack definitions as runtime rule and tool definitions.
type Registry struct {
	packs        []Definition
	capabilities []toolchain.Capability
	rules        []style.RuleDefinition
}

// Packs returns the packs registered in the registry.
func (registry Registry) Packs() (packs []Definition) {
	return CloneDefinitions(registry.packs)
}

// ToolCapabilities returns the tool capabilities registered in the registry.
func (registry Registry) ToolCapabilities() (capabilities []toolchain.Capability) {
	return append([]toolchain.Capability{}, registry.capabilities...)
}

// Rules returns the rule definitions registered in the registry.
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

func buildRegistry(packs []Definition) (registry Registry) {
	registry.packs = CloneDefinitions(packs)

	toolByID := make(map[string]toolchain.Capability)
	for _, pack := range packs {
		for _, tool := range pack.Tools {
			toolByID[tool.ID] = tool
		}

		registry.rules = append(registry.rules, CloneRules(pack.Rules)...)
	}

	toolIDs := make([]string, 0, len(toolByID))
	for toolID := range toolByID {
		toolIDs = append(toolIDs, toolID)
	}
	sort.Strings(toolIDs)

	registry.capabilities = make([]toolchain.Capability, 0, len(toolIDs))
	for _, toolID := range toolIDs {
		registry.capabilities = append(registry.capabilities, toolByID[toolID])
	}

	return registry
}

/* ----------------------------------------- Validation ----------------------------------------- */

func validateRegistry(registry Registry) (err error) {
	if err = validatePackFileSets(registry.packs); err != nil {
		return err
	}

	seenToolIDs := make(map[string]bool, len(registry.capabilities))
	for _, tool := range registry.capabilities {
		if tool.ID == "" {
			return fmt.Errorf("pack registry contains an empty tool id")
		}

		if tool.Name == "" {
			return fmt.Errorf("pack registry contains a tool %q with an empty name", tool.ID)
		}

		if seenToolIDs[tool.ID] {
			return fmt.Errorf("duplicate tool id %q", tool.ID)
		}

		seenToolIDs[tool.ID] = true
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
