package coverage

import "testing"

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
	plan := loadPlan(t)
	automated := ruleIDsByRequirement(plan.Rules)

	report := loadCoverageReport(t)
	for _, requirement := range report.Requirements {
		if requirement.Mode != ModeReviewOnly ||
			len(automated[requirement.ID]) == 0 {
			continue
		}

		t.Fatalf("review-only requirement %q is also marked automated", requirement.ID)
	}
}

func TestUnannotatedOutstandingRequirementsDefaultToDeferredManual(t *testing.T) {
	report := loadCoverageReport(t)
	for _, requirement := range report.Requirements {
		if requirement.Mode != ModeManualDeferred {
			continue
		}

		if requirement.Reason == "" {
			t.Fatalf("deferred manual requirement %q should include an explanation", requirement.ID)
		}
	}
}
