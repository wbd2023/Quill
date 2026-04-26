package report

import (
	"ciphera/tools/internal/coverage"
	"ciphera/tools/internal/styleguide"
)

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
