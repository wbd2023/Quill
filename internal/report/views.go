package report

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/coverage"
	"ciphera/tools/internal/styleguide"
)

/* ---------------------------------------- View Builders --------------------------------------- */

func NewCheckEntry(
	rule contract.Rule,
	status contract.CheckStatus,
	result contract.ExecutionResult,
) (entry CheckEntry) {
	return CheckEntry{
		Rule:   NewRuleSummary(rule),
		Status: status,
		Result: result,
	}
}

func NewRuleSummary(rule contract.Rule) (summary RuleSummary) {
	return RuleSummary{
		ID:             rule.ID,
		Name:           rule.Name,
		Group:          rule.Group,
		Level:          rule.Level,
		Scope:          rule.Scope,
		RequirementIDs: append([]string{}, rule.RequirementIDs...),
	}
}

func NewCheckView(result CheckResult) (view CheckView) {
	view = CheckView{
		Result:  result,
		Summary: result.Summary(),
		Groups:  make([]CheckGroup, 0),
	}

	for _, entry := range result.Entries {
		if len(view.Groups) == 0 ||
			entry.Rule.Group != view.Groups[len(view.Groups)-1].Group {
			view.Groups = append(view.Groups, CheckGroup{
				Group:   entry.Rule.Group,
				Entries: make([]CheckEntry, 0),
			})
		}

		lastIndex := len(view.Groups) - 1
		view.Groups[lastIndex].Entries = append(view.Groups[lastIndex].Entries, entry)
	}

	return view
}

func NewCoverageView(report coverage.Report) (view CoverageView) {
	view = CoverageView{
		Report:            report,
		Outstanding:       make([]coverage.Requirement, 0),
		OutstandingByMode: make(map[string]int),
	}

	for _, requirement := range report.Requirements {
		switch requirement.Mode {
		case styleguide.VerificationAutomated:
			view.RequirementTotals.Automated++
		case styleguide.VerificationReviewOnly:
			view.RequirementTotals.ReviewOnly++
		case styleguide.VerificationManualDeferred:
			view.RequirementTotals.ManualDeferred++
		}

		if requirement.Mode == styleguide.VerificationAutomated {
			continue
		}

		view.Outstanding = append(view.Outstanding, requirement)
		view.OutstandingByMode[string(requirement.Mode)]++
	}

	for _, section := range report.Sections {
		switch section.Status {
		case coverage.StatusAutomated:
			view.SectionTotals.Automated++

		case coverage.StatusPartial:
			view.SectionTotals.Partial++

		case coverage.StatusReviewOnly:
			view.SectionTotals.ReviewOnly++

		case coverage.StatusManual:
			view.SectionTotals.Manual++
		}
	}

	return view
}

func NewToolchainView(result ToolchainResult) (view ToolchainView) {
	view = ToolchainView{
		Result:   result,
		AllValid: true,
	}
	for _, status := range result.Statuses {
		if status.Valid {
			continue
		}

		view.AllValid = false
		break
	}

	return view
}
