package styleguide

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/rulepack"
)

func TestCoverageIncludesEveryStyleHeading(t *testing.T) {
	headings := extractHeadings(readStyleGuideForTest(t))
	covered := make(map[string]bool)

	report, err := Coverage(fixtures.RepoRoot(t))
	if err != nil {
		t.Fatalf("coverage: %v", err)
	}

	for _, entry := range report.Sections {
		covered[entry.Section] = true
	}

	for _, heading := range headings {
		section := heading.Section
		if covered[section] {
			continue
		}

		t.Fatalf("STYLE.md section %q missing from coverage index", section)
	}
}

func TestRequirementsCoverEveryStyleRequirement(t *testing.T) {
	documented := extractRequirements(readStyleGuideForTest(t))
	covered := make(map[string]bool)

	requirements, err := Requirements(fixtures.RepoRoot(t))
	if err != nil {
		t.Fatalf("requirements: %v", err)
	}

	for _, requirement := range requirements {
		covered[requirement.ID] = true
	}

	for _, requirement := range documented {
		if covered[requirement.ID] {
			continue
		}

		t.Fatalf("STYLE.md requirement %q missing from coverage model", requirement.ID)
	}
}

func TestReviewOnlyRequirementsAreNotAlsoAutomated(t *testing.T) {
	automated, err := ruleIDsByRequirement(fixtures.RepoRoot(t))
	if err != nil {
		t.Fatalf("ruleIDsByRequirement: %v", err)
	}

	requirements, err := Requirements(fixtures.RepoRoot(t))
	if err != nil {
		t.Fatalf("requirements: %v", err)
	}

	for _, requirement := range requirements {
		if requirement.Mode != VerificationReviewOnly || len(automated[requirement.ID]) == 0 {
			continue
		}

		t.Fatalf("review-only requirement %q is also marked automated", requirement.ID)
	}
}

func TestUnannotatedOutstandingRequirementsDefaultToDeferredManual(t *testing.T) {
	requirements, err := Requirements(fixtures.RepoRoot(t))
	if err != nil {
		t.Fatalf("requirements: %v", err)
	}

	for _, requirement := range requirements {
		if requirement.Mode != VerificationManualDeferred {
			continue
		}

		if requirement.Reason == "" {
			t.Fatalf("deferred manual requirement %q should include an explanation", requirement.ID)
		}
	}
}

/* --------------------------------------- Rule References -------------------------------------- */

func TestAutomatedRequirementsReferenceRules(t *testing.T) {
	ruleIDs := make(map[string]bool)
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
		ruleIDs[rule.ID] = true
	}

	requirements, err := Requirements(fixtures.RepoRoot(t))
	if err != nil {
		t.Fatalf("requirements: %v", err)
	}

	for _, requirement := range requirements {
		if requirement.Mode != VerificationAutomated {
			continue
		}

		if len(requirement.RuleIDs) > 0 {
			for _, ruleID := range requirement.RuleIDs {
				if ruleIDs[ruleID] {
					continue
				}

				t.Fatalf("requirement %q references unknown rule ID %q", requirement.ID, ruleID)
			}

			continue
		}

		t.Fatalf("automated requirement %q must reference at least one rule", requirement.ID)
	}
}
