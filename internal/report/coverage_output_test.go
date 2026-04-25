package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"ciphera/tools/internal/styleguide"
)

/* --------------------------------------- Coverage Output -------------------------------------- */

func TestWriteCoverageText(t *testing.T) {
	var buffer bytes.Buffer

	coverage := styleguide.CoverageReport{
		Requirements: []styleguide.Requirement{
			{
				ID:      "3.2.ctx-first",
				Section: "3.2",
				Text:    "`ctx context.Context` MUST be the first parameter when present.",
				Mode:    styleguide.VerificationAutomated,
				RuleIDs: []string{"go-policy"},
			},
			{
				ID:      "5.1.explain-security-plainly",
				Section: "5.1",
				Text:    "Security concepts SHOULD be explained in plain language.",
				Mode:    styleguide.VerificationReviewOnly,
				Reason:  "Plain-language quality is a writing judgement rather than a lint rule.",
			},
		},
		Sections: []styleguide.SectionCoverage{
			{
				Section:          "3.2",
				Title:            "Context, resources, and concurrency",
				Status:           styleguide.CoverageAutomated,
				RequirementCount: 1,
				AutomatedCount:   1,
			},
			{
				Section:          "5.1",
				Title:            "Audience",
				Status:           styleguide.CoverageReviewOnly,
				RequirementCount: 1,
				ReviewOnlyCount:  1,
			},
		},
	}

	if err := WriteCoverage(&buffer, FormatText, NewCoverageView(coverage), true); err != nil {
		t.Fatalf("WriteCoverage: %v", err)
	}

	output := buffer.String()
	if output != readGoldenOutput(t, "coverage.txt") {
		t.Fatalf("unexpected coverage output:\n%s", output)
	}
}

func TestWriteCoverageJSON(t *testing.T) {
	var buffer bytes.Buffer

	view := NewCoverageView(styleguide.CoverageReport{
		Requirements: []styleguide.Requirement{
			{
				ID:      "3.2.ctx-first",
				Section: "3.2",
				Mode:    styleguide.VerificationAutomated,
			},
		},
	})
	if err := WriteCoverage(&buffer, FormatJSON, view, false); err != nil {
		t.Fatalf("WriteCoverage: %v", err)
	}

	var envelope struct {
		Coverage CoverageView `json:"coverage"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode coverage json: %v", err)
	}

	if len(envelope.Coverage.Report.Requirements) != 1 {
		t.Fatalf("unexpected coverage payload: %+v", envelope.Coverage)
	}
}
