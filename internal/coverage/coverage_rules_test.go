package coverage

import (
	"slices"
	"testing"
)

func TestAutomatedRequirementsReferenceRules(t *testing.T) {
	ruleIDs := make(map[string]bool)
	plan := loadPlan(t)
	for _, rule := range plan.Rules {
		ruleIDs[rule.ID] = true
	}

	report := loadCoverageReport(t)
	for _, requirement := range report.Requirements {
		if requirement.Mode != ModeAutomated {
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

func TestQuillCoverageUsesGranularRuleBindings(t *testing.T) {
	report := loadCoverageReport(t)
	expected := map[string]string{
		"3.8.inline-comments-lowercase": "go/comments",
		"3.1.adapters-wrap-with-cause":  "go/errors",
		"3.4.explicit-parameter-types":  "go/parameters",
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
