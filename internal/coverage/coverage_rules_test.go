package coverage

import (
	"slices"
	"testing"
)

func TestAutomatedRequirementsReferenceRules(t *testing.T) {
	ruleIDs := make(map[string]bool)
	effective := loadEffectiveConfig(t)
	for _, rule := range effective.Rules {
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
