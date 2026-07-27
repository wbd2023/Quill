package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/wbd2023/quill/internal/coverage"
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
	if err := WriteCoverage(&buffer, "coverage", FormatText, view, true); err != nil {
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
	if err := WriteCoverage(&buffer, "coverage", FormatJSON, view, false); err != nil {
		t.Fatalf("WriteCoverage: %v", err)
	}

	var envelope struct {
		SchemaVersion int    `json:"schema_version"`
		Command       string `json:"command"`
		Status        string `json:"status"`
		Result        struct {
			Report struct {
				Requirements []struct {
					ID string `json:"id"`
				} `json:"requirements"`
			} `json:"report"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode coverage json: %v", err)
	}

	if envelope.SchemaVersion != SchemaVersion || envelope.Command != "coverage" ||
		envelope.Status != StatusOK {
		t.Fatalf("unexpected envelope header: %+v", envelope)
	}

	if len(envelope.Result.Report.Requirements) != 1 ||
		envelope.Result.Report.Requirements[0].ID != "3.2.ctx-first" {
		t.Fatalf("unexpected coverage payload: %+v", envelope.Result.Report)
	}
}
