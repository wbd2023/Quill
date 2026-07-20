package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/wbd2023/Quill/internal/coverage"
)

/* --------------------------------------- Coverage Output -------------------------------------- */

func TestWriteCoverageText(t *testing.T) {
	var buffer bytes.Buffer

	coverageReport := coverage.Report{
		Requirements: []coverage.Requirement{
			{
				ID:      "3.2.ctx-first",
				Section: "3.2",
				Text:    "`ctx context.Context` MUST be the first parameter when present.",
				Mode:    coverage.ModeAutomated,
				RuleIDs: []string{"go-policy"},
			},
			{
				ID:      "5.1.explain-security-plainly",
				Section: "5.1",
				Text:    "Security concepts SHOULD be explained in plain language.",
				Mode:    coverage.ModeReviewOnly,
				Reason:  "Plain-language quality is a writing judgement rather than a lint rule.",
			},
		},
		Sections: []coverage.Section{
			{
				Section:          "3.2",
				Title:            "Context, resources, and concurrency",
				Status:           coverage.StatusAutomated,
				RequirementCount: 1,
				AutomatedCount:   1,
			},
			{
				Section:          "5.1",
				Title:            "Audience",
				Status:           coverage.StatusReviewOnly,
				RequirementCount: 1,
				ReviewOnlyCount:  1,
			},
		},
	}

	view := NewCoverageView(coverageReport)
	if err := WriteCoverage(&buffer, FormatText, view, true); err != nil {
		t.Fatalf("WriteCoverage: %v", err)
	}

	output := buffer.String()
	if output != readGoldenOutput(t, "coverage.txt") {
		t.Fatalf("unexpected coverage output:\n%s", output)
	}
}

func TestWriteCoverageJSON(t *testing.T) {
	var buffer bytes.Buffer

	view := NewCoverageView(coverage.Report{
		Requirements: []coverage.Requirement{
			{
				ID:      "3.2.ctx-first",
				Section: "3.2",
				Mode:    coverage.ModeAutomated,
			},
		},
	})
	if err := WriteCoverage(&buffer, FormatJSON, view, false); err != nil {
		t.Fatalf("WriteCoverage: %v", err)
	}

	var envelope struct {
		Coverage struct {
			Report struct {
				Requirements []struct {
					ID string `json:"id"`
				} `json:"requirements"`
			} `json:"report"`
		} `json:"coverage"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode coverage json: %v", err)
	}

	if len(envelope.Coverage.Report.Requirements) != 1 ||
		envelope.Coverage.Report.Requirements[0].ID != "3.2.ctx-first" {
		t.Fatalf("unexpected coverage payload: %+v", envelope.Coverage)
	}
}
