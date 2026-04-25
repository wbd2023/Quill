package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"ciphera/tools/internal/contract"
)

/* ---------------------------------------- Check Output ---------------------------------------- */

func TestWriteCheckText(t *testing.T) {
	var buffer bytes.Buffer

	result := CheckResult{
		Entries: []CheckEntry{
			{
				Rule: contract.Rule{
					RuleDefinition: contract.RuleDefinition{
						ID:    "toolchain",
						Name:  "Pinned toolchain",
						Group: contract.RuleGroupControlPlane,
					},
					RequirementIDs: []string{"0.1.security-first"},
				},
				Status: CheckStatusPass,
			},
			{
				Rule: contract.Rule{
					RuleDefinition: contract.RuleDefinition{
						ID:    "markdown",
						Name:  "markdownlint",
						Group: contract.RuleGroupExternal,
					},
					RequirementIDs: []string{"5.2.concise-and-clear"},
				},
				Status: CheckStatusFail,
				Output: "missing from PATH",
			},
		},
	}

	summary, err := WriteCheck(&buffer, FormatText, NewCheckView(result), true)
	if err != nil {
		t.Fatalf("WriteCheck: %v", err)
	}

	if summary.Failed != 1 || summary.Passed != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	output := buffer.String()
	if output != readGoldenOutput(t, "check.txt") {
		t.Fatalf("unexpected check output:\n%s", output)
	}
}

func TestWriteCheckJSON(t *testing.T) {
	var buffer bytes.Buffer

	view := NewCheckView(CheckResult{
		Entries: []CheckEntry{
			{
				Rule: contract.Rule{
					RuleDefinition: contract.RuleDefinition{
						ID:    "toolchain",
						Name:  "Pinned toolchain",
						Group: contract.RuleGroupControlPlane,
					},
					RequirementIDs: []string{"0.1.security-first"},
				},
				Status: CheckStatusPass,
			},
		},
	})
	summary, err := WriteCheck(&buffer, FormatJSON, view, false)
	if err != nil {
		t.Fatalf("WriteCheck: %v", err)
	}

	if summary.Passed != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	var envelope struct {
		Check CheckView `json:"check"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode check json: %v", err)
	}

	if envelope.Check.Summary.Passed != 1 {
		t.Fatalf("unexpected JSON summary: %+v", envelope.Check.Summary)
	}
}
