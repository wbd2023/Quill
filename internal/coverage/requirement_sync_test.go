package coverage

import (
	"testing"

	"github.com/wbd2023/Quill/internal/pack"
	"github.com/wbd2023/Quill/internal/pack/shipped"
	"github.com/wbd2023/Quill/internal/profile"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/styleguide"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
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

	config, err = pack.ResolvePacks(config, registry.Packs())
	if err != nil {
		t.Fatalf("ResolvePacks: %v", err)
	}

	compiled, err := profile.Compile(config, registry.Definitions())
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
		id, err := style.ParseRequirementID(requirement.ID)
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
	document, err := styleguide.Load(testutil.RepositoryRoot(t), styleguide.Config{
		Filename: config.StyleGuide.Path,
	})
	if err != nil {
		t.Fatalf("load STYLE.md: %v", err)
	}

	return document.Requirements
}
