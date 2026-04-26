package coverage

import (
	"slices"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/styleguide"
)

func TestCoverageIncludesEveryStyleHeading(t *testing.T) {
	document := loadDocument(t)
	covered := make(map[string]bool)

	report := loadCoverageReport(t)
	for _, entry := range report.Sections {
		covered[entry.Section] = true
	}

	for _, heading := range document.Headings {
		section := heading.Section
		if covered[section] {
			continue
		}

		t.Fatalf("STYLE.md section %q missing from coverage index", section)
	}
}

func TestRequirementsCoverEveryStyleRequirement(t *testing.T) {
	document := loadDocument(t)
	covered := make(map[string]bool)

	report := loadCoverageReport(t)
	for _, requirement := range report.Requirements {
		covered[requirement.ID] = true
	}

	for _, requirement := range document.Requirements {
		if covered[requirement.ID] {
			continue
		}

		t.Fatalf("STYLE.md requirement %q missing from coverage model", requirement.ID)
	}
}

func TestReviewOnlyRequirementsAreNotAlsoAutomated(t *testing.T) {
	effective := loadEffectiveConfig(t)
	automated := ruleIDsByRequirement(effective.Rules)

	report := loadCoverageReport(t)
	for _, requirement := range report.Requirements {
		if requirement.Mode != styleguide.VerificationReviewOnly ||
			len(automated[requirement.ID]) == 0 {
			continue
		}

		t.Fatalf("review-only requirement %q is also marked automated", requirement.ID)
	}
}

func TestUnannotatedOutstandingRequirementsDefaultToDeferredManual(t *testing.T) {
	report := loadCoverageReport(t)
	for _, requirement := range report.Requirements {
		if requirement.Mode != styleguide.VerificationManualDeferred {
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
	effective := loadEffectiveConfig(t)
	for _, rule := range effective.Rules {
		ruleIDs[rule.ID] = true
	}

	report := loadCoverageReport(t)
	for _, requirement := range report.Requirements {
		if requirement.Mode != styleguide.VerificationAutomated {
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

func TestGoStyleCoverageUsesGranularRuleBindings(t *testing.T) {
	report := loadCoverageReport(t)
	expected := map[string]string{
		"2.2.structured-logs":            "go/logging",
		"3.8.constructor-category-order": "go/parameters",
		"1.8.blank-line-between-guards":  "go/guard-clause-spacing",
	}

	for requirementID, ruleID := range expected {
		requirement, found := coverageRequirementByID(report, requirementID)
		if !found {
			t.Fatalf("requirement %q missing from coverage report", requirementID)
		}

		if !slices.Contains(requirement.RuleIDs, ruleID) {
			t.Fatalf(
				"requirement %q covered by %v, want %q",
				requirementID,
				requirement.RuleIDs,
				ruleID,
			)
		}
	}
}

func TestGoStyleCoverageDoesNotReferenceLegacyMonolithRules(t *testing.T) {
	report := loadCoverageReport(t)
	for _, requirement := range report.Requirements {
		for _, ruleID := range requirement.RuleIDs {
			switch ruleID {
			case "go/style-app", "go/style-tools", "go/lint-app", "go/lint-tools":
				t.Fatalf("requirement %q references legacy rule %q", requirement.ID, ruleID)
			}
		}
	}
}

func loadDocument(t *testing.T) (document styleguide.Document) {
	t.Helper()

	config := profiles.Current(t)
	document, err := styleguide.Load(fixtures.RepoRoot(t), styleguide.Config{
		Path:                config.StyleGuide.Path,
		RequirementIDFormat: config.StyleGuide.RequirementIDFormat,
	})
	if err != nil {
		t.Fatalf("styleguide.Load: %v", err)
	}

	return document
}

func loadEffectiveConfig(t *testing.T) (effective contract.EffectiveConfig) {
	t.Helper()

	config := profiles.Current(t)
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err = profile.Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("profile.Compile: %v", err)
	}

	return effective
}

func loadCoverageReport(t *testing.T) (report Report) {
	t.Helper()

	return Build(loadDocument(t), loadEffectiveConfig(t).Rules)
}

func coverageRequirementByID(
	report Report,
	requirementID string,
) (requirement Requirement, found bool) {
	for _, requirement := range report.Requirements {
		if requirement.ID == requirementID {
			return requirement, true
		}
	}

	return Requirement{}, false
}
