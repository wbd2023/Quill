package report

import (
	"ciphera/tools/internal/coverage"
)

type CoverageTotals struct {
	Automated      int `json:"automated"`
	ReviewOnly     int `json:"review_only"`
	ManualDeferred int `json:"manual_deferred"`
}

type SectionTotals struct {
	Automated  int `json:"automated"`
	Partial    int `json:"partial"`
	ReviewOnly int `json:"review_only"`
	Manual     int `json:"manual"`
}

type CoverageView struct {
	Report            coverage.Report
	RequirementTotals CoverageTotals
	SectionTotals     SectionTotals
	Outstanding       []coverage.Requirement
	OutstandingByMode map[string]int
}

func NewCoverageView(report coverage.Report) (view CoverageView) {
	view = CoverageView{
		Report:            report,
		Outstanding:       make([]coverage.Requirement, 0),
		OutstandingByMode: make(map[string]int),
	}

	for _, requirement := range report.Requirements {
		switch requirement.Mode {
		case coverage.ModeAutomated:
			view.RequirementTotals.Automated++
		case coverage.ModeReviewOnly:
			view.RequirementTotals.ReviewOnly++
		case coverage.ModeManualDeferred:
			view.RequirementTotals.ManualDeferred++
		}

		if requirement.Mode == coverage.ModeAutomated {
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
