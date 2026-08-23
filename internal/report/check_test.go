package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/style"
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
				Result: style.ExecutionResult{
					Diagnostics: []style.Diagnostic{
						{
							File:    "docs/example.md",
							Range:   style.Range{Start: style.Position{Line: 5, Column: 3}},
							Message: "line too long",
						},
						{Message: "missing from PATH"},
					},
				},
			},
		},
	}

	summary, err := WriteCheck(
		&buffer, testEnvelopeMetadata("check"), FormatText, NewCheckView(result), true,
	)
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
	summary, err := WriteCheck(&buffer, testEnvelopeMetadata("check"), FormatJSON, view, false)
	if err != nil {
		t.Fatalf("WriteCheck: %v", err)
	}

	if summary.Passed != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	var envelope struct {
		SchemaVersion int    `json:"schema_version"`
		Command       string `json:"command"`
		Status        string `json:"status"`
		Result        struct {
			Summary CheckSummary `json:"summary"`
			Result  struct {
				Entries []struct {
					RuleID string `json:"rule_id"`
					Name   string `json:"name"`
				} `json:"entries"`
			} `json:"result"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode check json: %v", err)
	}

	if envelope.SchemaVersion != SchemaVersion {
		t.Fatalf("unexpected schema version: %d", envelope.SchemaVersion)
	}

	if envelope.Command != "check" || envelope.Status != StatusOK {
		t.Fatalf("unexpected envelope header: command=%q status=%q",
			envelope.Command, envelope.Status)
	}

	if envelope.Result.Summary.Passed != 1 {
		t.Fatalf("unexpected JSON summary: %+v", envelope.Result.Summary)
	}

	if len(envelope.Result.Result.Entries) != 1 ||
		envelope.Result.Result.Entries[0].RuleID != "toolchain" {
		t.Fatalf("unexpected JSON entries: %+v", envelope.Result.Result.Entries)
	}

	for _, forbidden := range []string{"spec", "fix_spec", "install_kind", "module_path"} {
		if strings.Contains(buffer.String(), forbidden) {
			t.Fatalf("check JSON leaked internal field %q: %s", forbidden, buffer.String())
		}
	}
}

func TestWriteCheckJSONEncodesRange(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	view := NewCheckView(CheckResult{
		Entries: []CheckEntry{
			{
				Rule: NewRuleSummary(style.Rule{
					ID:    "text",
					Name:  "line length",
					Group: style.RuleGroup("external_tools"),
				}),
				Status: style.CheckStatusFail,
				Result: style.ExecutionResult{
					Diagnostics: []style.Diagnostic{
						{
							Code:    "text/line-length/too-long",
							File:    "docs/example.md",
							Range:   style.Range{Start: style.Position{Line: 5, Column: 3}},
							Message: "line too long",
						},
						{Code: "text/line-length/too-long", Message: "missing from PATH"},
					},
				},
			},
		},
	})

	if _, err := WriteCheck(
		&buffer, testEnvelopeMetadata("check"), FormatJSON, view, false,
	); err != nil {
		t.Fatalf("WriteCheck: %v", err)
	}

	var envelope struct {
		Result struct {
			Groups []struct {
				Entries []struct {
					Diagnostics []map[string]json.RawMessage `json:"diagnostics"`
				} `json:"entries"`
			} `json:"groups"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode check json: %v", err)
	}

	if len(envelope.Result.Groups) != 1 || len(envelope.Result.Groups[0].Entries) != 1 {
		t.Fatalf("unexpected entry shape: %+v", envelope)
	}

	diagnostics := envelope.Result.Groups[0].Entries[0].Diagnostics
	if len(diagnostics) != 2 {
		t.Fatalf("expected two diagnostics, got %d", len(diagnostics))
	}

	// A known location renders a structured range: start carries line and column; end is omitted
	// when the extent is unknown.
	known, hasRange := diagnostics[0]["range"]
	if !hasRange {
		t.Fatalf("expected range for known-location diagnostic: %s", diagnostics[0])
	}

	var rangePayload struct {
		Start struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"start"`
		End *struct {
			Line int `json:"line"`
		} `json:"end"`
	}
	if err := json.Unmarshal(known, &rangePayload); err != nil {
		t.Fatalf("decode range: %v", err)
	}
	if rangePayload.Start.Line != 5 || rangePayload.Start.Column != 3 {
		t.Fatalf("unexpected range start: %+v", rangePayload.Start)
	}
	if rangePayload.End != nil {
		t.Fatalf("expected unknown range end to be omitted, got %+v", rangePayload.End)
	}

	// An unknown location omits the range key entirely rather than emitting zero values.
	if _, hasRange := diagnostics[1]["range"]; hasRange {
		t.Fatalf("expected no range for unknown-location diagnostic: %s", diagnostics[1])
	}
}

func TestWriteCheckJSONSerializesHelpURL(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	view := NewCheckView(CheckResult{
		Entries: []CheckEntry{
			{
				Rule: NewRuleSummary(style.Rule{
					ID:    "external",
					Name:  "external check",
					Group: style.RuleGroup("external"),
				}),
				Status: style.CheckStatusFail,
				Result: style.ExecutionResult{
					Diagnostics: []style.Diagnostic{{
						Code:    "external/rule",
						File:    "internal/service/users.go",
						Range:   style.Range{Start: style.Position{Line: 42, Column: 5}},
						Message: "Use the repository adapter.",
						HelpURL: "https://example.invalid/rules/external-rule",
					}},
				},
			},
		},
	})

	if _, err := WriteCheck(
		&buffer, testEnvelopeMetadata("check"), FormatJSON, view, false,
	); err != nil {
		t.Fatalf("WriteCheck: %v", err)
	}

	var envelope struct {
		Result struct {
			Groups []struct {
				Entries []struct {
					Diagnostics []struct {
						HelpURL string `json:"help_url"`
					} `json:"diagnostics"`
				} `json:"entries"`
			} `json:"groups"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode check json: %v", err)
	}

	diagnostics := envelope.Result.Groups[0].Entries[0].Diagnostics
	if len(diagnostics) != 1 {
		t.Fatalf("expected one diagnostic, got %d", len(diagnostics))
	}
	const want = "https://example.invalid/rules/external-rule"
	if diagnostics[0].HelpURL != want {
		t.Fatalf("expected help_url %q serialized, got %q", want, diagnostics[0].HelpURL)
	}
}
