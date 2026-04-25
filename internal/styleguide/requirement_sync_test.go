package styleguide

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/rulepack"
)

/* -------------------------------------------- Sync -------------------------------------------- */

func TestRuleRequirementIDsExistInStyleGuide(t *testing.T) {
	requirements := make(map[string]bool)
	for _, requirement := range extractRequirements(readStyleGuideForTest(t)) {
		requirements[requirement.ID] = true
	}

	policy := profiles.Current(t)

	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err := policy.Compile(registry)
	if err != nil {
		t.Fatalf("profile.Compile: %v", err)
	}

	for _, rule := range effective.Rules {
		for _, requirementID := range rule.RequirementIDs {
			if requirements[requirementID] {
				continue
			}

			t.Fatalf("rule %q references missing STYLE.md requirement %q", rule.ID, requirementID)
		}
	}
}

func TestStyleGuideRequirementIDsMatchTheirSection(t *testing.T) {
	for _, requirement := range extractRequirements(readStyleGuideForTest(t)) {
		expectedSection := RequirementSection(requirement.ID)
		if expectedSection == requirement.Section {
			continue
		}

		t.Fatalf(
			"requirement %q appears under STYLE.md section %q, expected %q",
			requirement.ID,
			requirement.Section,
			expectedSection,
		)
	}
}

func TestStyleGuideNonAutomatedRequirementsHaveReasons(t *testing.T) {
	for _, requirement := range extractRequirements(readStyleGuideForTest(t)) {
		if requirement.Mode == "" {
			continue
		}
		if requirement.Reason != "" {
			continue
		}

		t.Fatalf("non-automated requirement %q must declare a reason", requirement.ID)
	}
}

func TestStyleGuideRequirementIDsAreUnique(t *testing.T) {
	seen := make(map[string]bool)

	for _, requirement := range extractRequirements(readStyleGuideForTest(t)) {
		if seen[requirement.ID] {
			t.Fatalf("duplicate STYLE.md requirement ID %q", requirement.ID)
		}

		seen[requirement.ID] = true
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func readStyleGuideForTest(t *testing.T) (contents string) {
	t.Helper()

	_, data, err := readStyleGuide(fixtures.RepoRoot(t))
	if err != nil {
		t.Fatalf("read STYLE.md: %v", err)
	}

	return string(data)
}
