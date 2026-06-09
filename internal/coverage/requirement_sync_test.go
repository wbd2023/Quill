package coverage

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/requirementid"
	"ciphera/tools/internal/styleguide"
)

/* -------------------------------------------- Sync -------------------------------------------- */

func TestRuleRequirementIDsExistInStyleGuide(t *testing.T) {
	requirements := make(map[string]bool)
	for _, requirement := range loadStyleRequirements(t) {
		requirements[requirement.ID] = true
	}

	config := profiles.Current(t)

	registry, err := shipped.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	compiled, err := profile.Compile(config, registry)
	if err != nil {
		t.Fatalf("profile.Compile: %v", err)
	}

	for _, rule := range compiled.Effective.Rules {
		for _, requirementID := range rule.RequirementIDs {
			if requirements[requirementID] {
				continue
			}

			t.Fatalf("rule %q references missing STYLE.md requirement %q", rule.ID, requirementID)
		}
	}
}

func TestStyleGuideRequirementIDsMatchTheirSection(t *testing.T) {
	for _, requirement := range loadStyleRequirements(t) {
		id, err := requirementid.Parse(requirement.ID, requirementid.SectionSlug)
		if err != nil {
			t.Fatalf("parse requirement ID %q: %v", requirement.ID, err)
		}

		if id.Section() == requirement.Section {
			continue
		}

		t.Fatalf(
			"requirement %q appears under STYLE.md section %q, expected %q",
			requirement.ID,
			requirement.Section,
			id.Section(),
		)
	}
}

func TestStyleGuideNonAutomatedRequirementsHaveReasons(t *testing.T) {
	for _, requirement := range loadStyleRequirements(t) {
		if !requirement.Review.Only {
			continue
		}
		if requirement.Review.Reason != "" {
			continue
		}

		t.Fatalf("non-automated requirement %q must declare a reason", requirement.ID)
	}
}

func TestStyleGuideRequirementIDsAreUnique(t *testing.T) {
	seen := make(map[string]bool)

	for _, requirement := range loadStyleRequirements(t) {
		if seen[requirement.ID] {
			t.Fatalf("duplicate STYLE.md requirement ID %q", requirement.ID)
		}

		seen[requirement.ID] = true
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func loadStyleRequirements(t *testing.T) (requirements []styleguide.Requirement) {
	t.Helper()

	config := profiles.Current(t)
	document, err := styleguide.Load(fixtures.RepositoryRoot(t), styleguide.Config{
		Filename: config.StyleGuide.Path,
		IDScheme: config.StyleGuide.IDScheme,
	})
	if err != nil {
		t.Fatalf("load STYLE.md: %v", err)
	}

	return document.Requirements
}
