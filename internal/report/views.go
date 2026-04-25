package report

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/styleguide"
)

/* ---------------------------------------- View Builders --------------------------------------- */

func NewCheckView(result CheckResult) (view CheckView) {
	view = CheckView{
		Result:  result,
		Summary: result.Summary(),
		Groups:  make([]CheckGroup, 0),
	}

	currentGroup := contract.RuleGroup("")
	for _, entry := range result.Entries {
		if len(view.Groups) == 0 || entry.Rule.Group != currentGroup {
			currentGroup = entry.Rule.Group
			view.Groups = append(view.Groups, CheckGroup{
				Group:   currentGroup,
				Entries: make([]CheckEntry, 0),
			})
		}

		lastIndex := len(view.Groups) - 1
		view.Groups[lastIndex].Entries = append(view.Groups[lastIndex].Entries, entry)
	}

	return view
}

func NewCoverageView(report styleguide.CoverageReport) (view CoverageView) {
	view = CoverageView{
		Report:            report,
		Outstanding:       make([]styleguide.Requirement, 0),
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
		case styleguide.CoverageAutomated:
			view.SectionTotals.Automated++

		case styleguide.CoveragePartial:
			view.SectionTotals.Partial++

		case styleguide.CoverageReviewOnly:
			view.SectionTotals.ReviewOnly++

		case styleguide.CoverageManual:
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
