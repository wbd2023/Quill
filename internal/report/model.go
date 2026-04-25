package report

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/styleguide"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	CheckStatusPass CheckStatus = "pass"
	CheckStatusWarn CheckStatus = "warn"
	CheckStatusFail CheckStatus = "fail"
	CheckStatusSkip CheckStatus = "skip"
)

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
)

/* -------------------------------------------- Types ------------------------------------------- */

type CheckStatus string

type OutputFormat string

type CheckEntry struct {
	Rule   contract.Rule
	Status CheckStatus
	Output string
}

type CheckResult struct {
	Entries []CheckEntry
}

type CheckSummary struct {
	Passed  int
	Warned  int
	Failed  int
	Skipped int
}

type CheckGroup struct {
	Group   contract.RuleGroup `json:"group"`
	Entries []CheckEntry       `json:"entries"`
}

type CheckView struct {
	Result  CheckResult  `json:"result"`
	Summary CheckSummary `json:"summary"`
	Groups  []CheckGroup `json:"groups"`
}

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
	Report            styleguide.CoverageReport `json:"report"`
	RequirementTotals CoverageTotals            `json:"requirement_totals"`
	SectionTotals     SectionTotals             `json:"section_totals"`
	Outstanding       []styleguide.Requirement  `json:"outstanding"`
	OutstandingByMode map[string]int            `json:"outstanding_by_mode"`
}

type ToolchainResult struct {
	Statuses []runtime.ToolStatus
}

type ToolchainView struct {
	Result   ToolchainResult `json:"result"`
	AllValid bool            `json:"all_valid"`
}

/* --------------------------------------- Check Summaries -------------------------------------- */

func (result CheckResult) Summary() (summary CheckSummary) {
	for _, entry := range result.Entries {
		switch entry.Status {
		case CheckStatusPass:
			summary.Passed++

		case CheckStatusWarn:
			summary.Warned++

		case CheckStatusFail:
			summary.Failed++

		case CheckStatusSkip:
			summary.Skipped++
		}
	}

	return summary
}
