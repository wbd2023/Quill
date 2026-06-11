package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"ciphera/tools/internal/style"
)

/* ---------------------------------------- Check Output ---------------------------------------- */

func TestWriteCheckText(t *testing.T) {
	var buffer bytes.Buffer

	result := CheckResult{
		Entries: []CheckEntry{
			{
				Rule: NewRuleSummary(style.Rule{
					ID:             "toolchain",
					Name:           "Pinned toolchain",
					Group:          style.RuleGroup("project"),
					RequirementIDs: []string{"0.1.security-first"},
				}),
				Status: style.CheckStatusPass,
			},
			{
				Rule: NewRuleSummary(style.Rule{
					ID:             "markdown",
					Name:           "markdownlint",
					Group:          style.RuleGroup("external_tools"),
					RequirementIDs: []string{"5.2.concise-and-clear"},
				}),
				Status: style.CheckStatusFail,
				Result: style.ExecutionResult{Output: "missing from PATH"},
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
				Rule: NewRuleSummary(style.Rule{
					ID:             "toolchain",
					Name:           "Pinned toolchain",
					Group:          style.RuleGroup("project"),
					RequirementIDs: []string{"0.1.security-first"},
				}),
				Status: style.CheckStatusPass,
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
		Check struct {
			Summary CheckSummary `json:"summary"`
			Result  struct {
				Entries []struct {
					RuleID string `json:"rule_id"`
					Name   string `json:"name"`
				} `json:"entries"`
			} `json:"result"`
		} `json:"check"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode check json: %v", err)
	}

	if envelope.Check.Summary.Passed != 1 {
		t.Fatalf("unexpected JSON summary: %+v", envelope.Check.Summary)
	}

	if len(envelope.Check.Result.Entries) != 1 ||
		envelope.Check.Result.Entries[0].RuleID != "toolchain" {
		t.Fatalf("unexpected JSON entries: %+v", envelope.Check.Result.Entries)
	}

	for _, forbidden := range []string{"spec", "fix_spec", "install_kind", "module_path"} {
		if strings.Contains(buffer.String(), forbidden) {
			t.Fatalf("check JSON leaked internal field %q: %s", forbidden, buffer.String())
		}
	}
}
