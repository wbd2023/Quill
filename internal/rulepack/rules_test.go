package rulepack_test

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/rulepack"
)

/* --------------------------------------- Rule Contracts --------------------------------------- */

func TestRegisteredRulesHaveUniqueIDs(t *testing.T) {
	registry, err := rulepack.DefaultRegistry(nil)
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
	registry, err := rulepack.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, rule := range registry.Rules() {
		for _, toolID := range rule.ToolIDs() {
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
	policy := profiles.Current(t)

	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err := policy.Compile(registry)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	if len(effective.Rules) != len(registry.Rules()) {
		t.Fatalf("expected %d active rules, got %d", len(registry.Rules()), len(effective.Rules))
	}

	for _, rule := range effective.Rules {
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
	registry, err := rulepack.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	validExecutors := map[string]bool{
		contract.ExecutorToolchain:      true,
		contract.ExecutorControlPlane:   true,
		contract.ExecutorFileCommand:    true,
		contract.ExecutorGoFormat:       true,
		contract.ExecutorGolangci:       true,
		contract.ExecutorGoStyle:        true,
		contract.ExecutorRepositoryScan: true,
	}

	for _, rule := range registry.Rules() {
		if !validExecutors[rule.Spec.Executor] {
			t.Fatalf("rule %q uses unsupported executor %q", rule.ID, rule.Spec.Executor)
		}

		if rule.FixSpec.Executor == "" || validExecutors[rule.FixSpec.Executor] {
			continue
		}

		t.Fatalf("rule %q uses unsupported fix executor %q", rule.ID, rule.FixSpec.Executor)
	}
}

func TestRuleGroupsRemainStable(t *testing.T) {
	registry, err := rulepack.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	seenGroups := make(map[contract.RuleGroup]bool)
	for _, rule := range registry.Rules() {
		seenGroups[rule.Group] = true
	}

	for _, group := range []contract.RuleGroup{
		contract.RuleGroupControlPlane,
		contract.RuleGroupExternal,
		contract.RuleGroupLanguage,
		contract.RuleGroupRepository,
	} {
		if seenGroups[group] {
			continue
		}

		t.Fatalf("missing rule group %q", group)
	}
}
