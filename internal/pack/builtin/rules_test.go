package builtin

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/profile"
)

/* --------------------------------------- Rule Contracts --------------------------------------- */

func TestRegisteredRulesHaveUniqueIDs(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	seenIDs := make(map[string]bool)
	for _, rule := range registry.Rules() {
		if seenIDs[rule.ID] {
			t.Fatalf("duplicate rule ID: %s", rule.ID)
		}

		seenIDs[rule.ID] = true
	}
}

func TestRegisteredRulesReferenceKnownTools(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, rule := range registry.Rules() {
		for _, toolID := range rule.CheckToolIDs() {
			if _, found := registry.ToolByID(toolID); found {
				continue
			}

			t.Fatalf("rule %q references unknown tool %q", rule.ID, toolID)
		}

		for _, toolID := range rule.FixToolIDs() {
			if _, found := registry.ToolByID(toolID); found {
				continue
			}

			t.Fatalf("rule %q fix references unknown tool %q", rule.ID, toolID)
		}
	}
}

func TestCurrentProfileBindsEveryRegisteredRule(t *testing.T) {
	config := profiles.Current(t)

	registry, err := DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	compiled, err := profile.Compile(config, registry)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	if len(compiled.Effective.Rules) != len(registry.Rules()) {
		t.Fatalf(
			"expected %d active rules, got %d",
			len(registry.Rules()),
			len(compiled.Effective.Rules),
		)
	}

	for _, rule := range compiled.Effective.Rules {
		if len(rule.RequirementIDs) == 0 {
			t.Fatalf("rule %q must reference at least one STYLE.md requirement", rule.ID)
		}

		seenRequirements := make(map[string]bool, len(rule.RequirementIDs))
		for _, requirementID := range rule.RequirementIDs {
			if seenRequirements[requirementID] {
				t.Fatalf("rule %q duplicates requirement %q", rule.ID, requirementID)
			}

			seenRequirements[requirementID] = true
		}
	}
}

func TestRegisteredRulesUseExpectedExecutors(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	validExecutors := map[contract.ExecutorKind]bool{
		contract.ExecutorToolchain:      true,
		contract.ExecutorProject:        true,
		contract.ExecutorFileCommand:    true,
		contract.ExecutorTargetCommand:  true,
		contract.ExecutorTargetCheck:    true,
		contract.ExecutorRepositoryScan: true,
	}

	for _, rule := range registry.Rules() {
		if !validExecutors[rule.Check.Kind] {
			t.Fatalf("rule %q uses unsupported executor %q", rule.ID, rule.Check.Kind)
		}

		if rule.Fix.Empty() || validExecutors[rule.Fix.Kind] {
			continue
		}

		t.Fatalf("rule %q uses unsupported fix executor %q", rule.ID, rule.Fix.Kind)
	}
}

func TestRuleGroupsRemainStable(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	seenGroups := make(map[contract.RuleGroup]bool)
	for _, rule := range registry.Rules() {
		seenGroups[rule.Group] = true
	}

	expectedGroups := []contract.RuleGroup{
		"project",
		"external_tools",
		"language",
		"text_scanners",
		"security_scanners",
		"vocabulary_scanners",
	}

	for _, group := range expectedGroups {
		if seenGroups[group] {
			continue
		}

		t.Fatalf("missing rule group %q", group)
	}
}
