package report

import (
	"io"

	"github.com/wbd2023/Quill/internal/coverage"
)

/* ------------------------------------------ JSON DTOs ----------------------------------------- */

type coverageJSON struct {
	Report            coverageReportJSON `json:"report"`
	RequirementTotals CoverageTotals     `json:"requirement_totals"`
	SectionTotals     SectionTotals      `json:"section_totals"`
	Outstanding       []requirementJSON  `json:"outstanding"`
	OutstandingByMode map[string]int     `json:"outstanding_by_mode"`
}

type coverageReportJSON struct {
	Requirements []requirementJSON `json:"requirements"`
	Sections     []sectionJSON     `json:"sections"`
}

type requirementJSON struct {
	ID      string   `json:"id"`
	Section string   `json:"section"`
	Text    string   `json:"text"`
	Mode    string   `json:"mode"`
	Reason  string   `json:"reason,omitempty"`
	RuleIDs []string `json:"rule_ids"`
}

type sectionJSON struct {
	Section             string          `json:"section"`
	Title               string          `json:"title"`
	Status              coverage.Status `json:"status"`
	RequirementCount    int             `json:"requirement_count"`
	AutomatedCount      int             `json:"automated_count"`
	ReviewOnlyCount     int             `json:"review_only_count"`
	ManualDeferredCount int             `json:"manual_deferred_count"`
}

/* ------------------------------------------ Rendering ----------------------------------------- */

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
