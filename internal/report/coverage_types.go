package report

import "ciphera/tools/internal/coverage"

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
