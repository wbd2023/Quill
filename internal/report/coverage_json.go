package report

import (
	"io"

	"ciphera/tools/internal/coverage"
)

func writeCoverageJSON(writer io.Writer, view CoverageView) (err error) {
	return writeJSON(writer, struct {
		Coverage coverageJSON `json:"coverage"`
	}{Coverage: newCoverageJSON(view)})
}

func newCoverageJSON(view CoverageView) (payload coverageJSON) {
	return coverageJSON{
		Report: coverageReportJSON{
			Requirements: requirementListJSON(view.Report.Requirements),
			Sections:     sectionListJSON(view.Report.Sections),
		},
		RequirementTotals: view.RequirementTotals,
		SectionTotals:     view.SectionTotals,
		Outstanding:       requirementListJSON(view.Outstanding),
		OutstandingByMode: cloneIntMap(view.OutstandingByMode),
	}
}

func requirementListJSON(requirements []coverage.Requirement) (payload []requirementJSON) {
	payload = make([]requirementJSON, 0, len(requirements))
	for _, requirement := range requirements {
		payload = append(payload, requirementJSON{
			ID:      requirement.ID,
			Section: requirement.Section,
			Text:    requirement.Text,
			Mode:    string(requirement.Mode),
			Reason:  requirement.Reason,
			RuleIDs: append([]string{}, requirement.RuleIDs...),
		})
	}

	return payload
}

func sectionListJSON(sections []coverage.Section) (payload []sectionJSON) {
	payload = make([]sectionJSON, 0, len(sections))
	for _, section := range sections {
		payload = append(payload, sectionJSON{
			Section:             section.Section,
			Title:               section.Title,
			Status:              section.Status,
			RequirementCount:    section.RequirementCount,
			AutomatedCount:      section.AutomatedCount,
			ReviewOnlyCount:     section.ReviewOnlyCount,
			ManualDeferredCount: section.ManualDeferredCount,
		})
	}

	return payload
}

func cloneIntMap(source map[string]int) (target map[string]int) {
	target = make(map[string]int, len(source))
	for key, value := range source {
		target[key] = value
	}

	return target
}
