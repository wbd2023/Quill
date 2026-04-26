package rulepack

import (
	"fmt"
	"sort"

	"ciphera/tools/internal/toolchain"
)

/* -------------------------------------- Registry Loading -------------------------------------- */

func DefaultRegistry(enabled []string) (registry Registry, err error) {
	packs := []Pack{
		controlPack(),
		textPack(),
		markdownPack(),
		shellPack(),
		goPack(),
		securityPack(),
		namingPack(),
	}

	if len(enabled) > 0 {
		packs, err = selectPacks(packs, enabled)
		if err != nil {
			return Registry{}, err
		}
	}

	registry = buildRegistry(packs)
	if err = validateRegistry(registry); err != nil {
		return Registry{}, err
	}

	return registry, nil
}

func selectPacks(available []Pack, enabled []string) (selected []Pack, err error) {
	packByID := make(map[string]Pack, len(available))
	for _, pack := range available {
		packByID[pack.ID] = pack
	}

	selected = make([]Pack, 0, len(enabled))
	seen := make(map[string]bool, len(enabled))
	for _, packID := range enabled {
		if seen[packID] {
			return nil, fmt.Errorf("duplicate rule pack %q", packID)
		}

		pack, found := packByID[packID]
		if !found {
			return nil, fmt.Errorf("unknown rule pack %q", packID)
		}

		seen[packID] = true
		selected = append(selected, pack)
	}

	return selected, nil
}

func buildRegistry(packs []Pack) (registry Registry) {
	registry.packs = append([]Pack{}, packs...)

	toolByID := make(map[string]toolchain.Capability)
	for _, pack := range packs {
		for _, tool := range pack.Tools {
			toolByID[tool.ID] = tool
		}

		registry.rules = append(registry.rules, pack.Rules...)
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

/* ------------------------------------- Registry Validation ------------------------------------ */

func validateRegistry(registry Registry) (err error) {
	if err = validatePackToolDefinitions(registry.packs); err != nil {
		return err
	}

	seenToolIDs := make(map[string]bool, len(registry.capabilities))
	for _, tool := range registry.capabilities {
		if tool.ID == "" {
			return fmt.Errorf("rule-pack registry contains an empty tool id")
		}

		if seenToolIDs[tool.ID] {
			return fmt.Errorf("duplicate tool id %q", tool.ID)
		}

		seenToolIDs[tool.ID] = true
	}

	seenRuleIDs := make(map[string]bool, len(registry.rules))
	for _, rule := range registry.rules {
		if rule.ID == "" {
			return fmt.Errorf("rule-pack registry contains an empty rule id")
		}

		if rule.Spec.Empty() {
			return fmt.Errorf("rule %q has no executor", rule.ID)
		}

		if seenRuleIDs[rule.ID] {
			return fmt.Errorf("duplicate rule id %q", rule.ID)
		}

		seenRuleIDs[rule.ID] = true
	}

	return nil
}

func validatePackToolDefinitions(packs []Pack) (err error) {
	toolByID := make(map[string]toolchain.Capability)
	for _, pack := range packs {
		for _, tool := range pack.Tools {
			existing, found := toolByID[tool.ID]
			if !found {
				toolByID[tool.ID] = tool
				continue
			}

			if existing != tool {
				return fmt.Errorf("tool %q has conflicting definitions", tool.ID)
			}
		}
	}

	return nil
}
