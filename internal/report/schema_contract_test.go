package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/wbd2023/quill/internal/coverage"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

/* ------------------------------------- Schema Contract --------------------------------------- */

// These tests defend the immutable JSON v1 result schemas documented in docs/cli-protocol.md. Each
// test renders a populated result through the real report renderer and asserts every documented
// required field is present with its documented JSON type.
//
// They never reject additive fields: a field is checked only for presence and type, never for the
// absence of extras. Additive extensions permitted within the same protocol version keep passing,
// matching the compatibility rule. A field dropped or retyped by the renderer fails immediately.

// jsonKind reports the JSON type name encoding/json produces when decoding into any.
func jsonKind(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case bool:
		return "bool"
	case float64:
		return "number"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return fmt.Sprintf("%T", value)
	}
}

// requireField asserts obj carries key with the documented JSON type.
func requireField(t *testing.T, obj map[string]any, key string, want string) {
	t.Helper()
	raw, ok := obj[key]
	if !ok {
		t.Fatalf("result missing documented field %q", key)
	}
	if got := jsonKind(raw); got != want {
		t.Fatalf("field %q is %s, want %s (value %v)", key, got, want, raw)
	}
}

// requireObject asserts key is a JSON object and returns it decoded.
func requireObject(t *testing.T, obj map[string]any, key string) map[string]any {
	t.Helper()
	requireField(t, obj, key, "object")
	return obj[key].(map[string]any)
}

// requireArray asserts key is a JSON array and returns it decoded.
func requireArray(t *testing.T, obj map[string]any, key string) []any {
	t.Helper()
	requireField(t, obj, key, "array")
	return obj[key].([]any)
}

// decodeResult decodes the result object from one rendered envelope.
func decodeResult(t *testing.T, output []byte) map[string]any {
	t.Helper()
	var envelope struct {
		Result map[string]any `json:"result"`
	}
	if err := json.Unmarshal(output, &envelope); err != nil {
		t.Fatalf("decode result envelope: %v\n%s", err, output)
	}
	if envelope.Result == nil {
		t.Fatalf("envelope missing result:\n%s", output)
	}
	return envelope.Result
}

func TestCheckResultSchemaContract(t *testing.T) {
	t.Parallel()

	view := NewCheckView(CheckResult{
		Entries: []CheckEntry{
			{
				Rule: NewRuleSummary(style.Rule{
					ID:             "gofmt",
					Name:           "gofmt",
					Group:          style.RuleGroup("formatting"),
					Enforcement:    style.EnforcementRequired,
					Scope:          style.Scope("all"),
					RequirementIDs: []string{"1.1.format"},
				}),
				Status: style.CheckStatusFail,
				Result: style.ExecutionResult{
					ExitCode: 2,
					Diagnostics: []style.Diagnostic{{
						Code:    "gofmt/diff",
						File:    "main.go",
						Range:   style.Range{Start: style.Position{Line: 1, Column: 1}},
						Message: "file is not gofmt-ed",
						HelpURL: "https://example.invalid/gofmt",
					}},
				},
				ExecutionError: errors.New("gofmt unavailable"),
			},
		},
	})

	var buffer bytes.Buffer
	if _, err := WriteCheck(&buffer, testEnvelopeMetadata("check"), FormatJSON, view, false); err != nil {
		t.Fatalf("WriteCheck: %v", err)
	}
	payload := decodeResult(t, buffer.Bytes())

	result := requireObject(t, payload, "result")

	requireField(t, result, "entries", "array")
	requireField(t, payload, "summary", "object")
	requireField(t, payload, "groups", "array")

	// CheckSummary carries capitalized field names; the doc lists them exactly.
	summary := requireObject(t, payload, "summary")
	for _, field := range []string{"Passed", "Warned", "Failed", "Blocked", "Skipped", "Errored"} {
		requireField(t, summary, field, "number")
	}

	entry := requireArray(t, result, "entries")[0].(map[string]any)
	for _, field := range []string{"rule_id", "name", "group", "enforcement", "scope", "status", "execution_error"} {
		requireField(t, entry, field, "string")
	}
	requireField(t, entry, "requirements", "array")
	requireField(t, entry, "diagnostics", "array")
	requireField(t, entry, "command", "object")

	command := requireObject(t, entry, "command")
	requireField(t, command, "exit_code", "number")
	requireField(t, command, "timed_out", "bool")
	requireField(t, command, "truncated", "bool")

	diagnostic := requireArray(t, entry, "diagnostics")[0].(map[string]any)
	for _, field := range []string{"code", "file", "message", "help_url"} {
		requireField(t, diagnostic, field, "string")
	}
	requireField(t, diagnostic, "range", "object")

	group := requireArray(t, payload, "groups")[0].(map[string]any)
	requireField(t, group, "group", "string")
	requireField(t, group, "entries", "array")
}

func TestCoverageResultSchemaContract(t *testing.T) {
	t.Parallel()

	view := NewCoverageView(coverage.Report{
		Requirements: []coverage.Requirement{{
			ID:      "1.1.format",
			Section: "1.1",
			Text:    "Files MUST be gofmt-ed.",
			Mode:    coverage.ModeAutomated,
			RuleIDs: []string{"gofmt"},
		}},
		Sections: []coverage.Section{{
			Section:             "1.1",
			Title:               "Formatting",
			Status:              coverage.StatusAutomated,
			RequirementCount:    1,
			AutomatedCount:      1,
			ReviewOnlyCount:     0,
			ManualDeferredCount: 0,
		}},
	})

	var buffer bytes.Buffer
	if err := WriteCoverage(&buffer, testEnvelopeMetadata("coverage"), FormatJSON, view, false); err != nil {
		t.Fatalf("WriteCoverage: %v", err)
	}
	result := decodeResult(t, buffer.Bytes())

	report := requireObject(t, result, "report")
	requireField(t, report, "requirements", "array")
	requireField(t, report, "sections", "array")

	requirement := requireArray(t, report, "requirements")[0].(map[string]any)
	for _, field := range []string{"id", "section", "text", "mode"} {
		requireField(t, requirement, field, "string")
	}
	requireField(t, requirement, "rule_ids", "array")

	section := requireArray(t, report, "sections")[0].(map[string]any)
	for _, field := range []string{"section", "title", "status"} {
		requireField(t, section, field, "string")
	}
	for _, field := range []string{
		"requirement_count", "automated_count", "review_only_count", "manual_deferred_count",
	} {
		requireField(t, section, field, "number")
	}

	requirementTotals := requireObject(t, result, "requirement_totals")
	for _, field := range []string{"automated", "review_only", "manual_deferred"} {
		requireField(t, requirementTotals, field, "number")
	}
	sectionTotals := requireObject(t, result, "section_totals")
	for _, field := range []string{"automated", "partial", "review_only", "manual"} {
		requireField(t, sectionTotals, field, "number")
	}
	requireField(t, result, "outstanding", "array")
	requireField(t, result, "outstanding_by_mode", "object")
}

func TestFixResultSchemaContract(t *testing.T) {
	t.Parallel()

	view := NewFixView(
		style.Scope("all"),
		false,
		[]toolchain.Status{{
			Tool:    toolchain.Tool{ID: "gofmt", Name: "gofmt", PinnedVersion: "1.24.5"},
			Path:    "/tools/gofmt",
			Version: "1.24.5",
			Valid:   false,
			Issue:   "not found",
		}},
		[]FixEntry{{
			Rule: NewRuleSummary(style.Rule{
				ID:          "gofmt",
				Name:        "gofmt",
				Group:       style.RuleGroup("formatting"),
				Enforcement: style.EnforcementRequired,
				Scope:       style.Scope("all"),
			}),
			Execution:      style.ExecutionResult{ExitCode: 2, TimedOut: true},
			ExecutionError: errors.New("gofmt unavailable"),
		}},
	)

	var buffer bytes.Buffer
	if _, err := WriteFix(&buffer, testEnvelopeMetadata("fix"), FormatJSON, view); err != nil {
		t.Fatalf("WriteFix: %v", err)
	}
	result := decodeResult(t, buffer.Bytes())

	requireField(t, result, "scope", "string")
	toolchainObject := requireObject(t, result, "toolchain")
	requireField(t, toolchainObject, "all_valid", "bool")
	requireField(t, toolchainObject, "statuses", "array")
	requireField(t, result, "rules", "array")

	status := requireArray(t, toolchainObject, "statuses")[0].(map[string]any)
	for _, field := range []string{"id", "name", "path", "version", "pinned_version", "issue"} {
		requireField(t, status, field, "string")
	}
	requireField(t, status, "valid", "bool")

	rule := requireArray(t, result, "rules")[0].(map[string]any)
	for _, field := range []string{"rule_id", "name", "group", "enforcement", "scope", "execution_error"} {
		requireField(t, rule, field, "string")
	}
	requireField(t, rule, "exit_code", "number")
	requireField(t, rule, "timed_out", "bool")
	requireField(t, rule, "truncated", "bool")
}

// TestToolchainResultSchemaContract covers doctor and install, which share the toolchain result.
func TestToolchainResultSchemaContract(t *testing.T) {
	t.Parallel()

	toolchainResult := ToolchainResult{
		AllValid: false,
		Statuses: []toolchain.Status{{
			Tool:    toolchain.Tool{ID: "gofmt", Name: "gofmt", PinnedVersion: "1.24.5"},
			Path:    "/tools/gofmt",
			Version: "1.24.5",
			Valid:   false,
			Issue:   "not found",
		}},
	}

	var buffer bytes.Buffer
	if _, err := WriteToolchain(&buffer, testEnvelopeMetadata("doctor"), FormatJSON, toolchainResult); err != nil {
		t.Fatalf("WriteToolchain: %v", err)
	}
	payload := decodeResult(t, buffer.Bytes())
	result := requireObject(t, payload, "result")

	requireField(t, result, "statuses", "array")
	requireField(t, payload, "all_valid", "bool")

	status := requireArray(t, result, "statuses")[0].(map[string]any)
	for _, field := range []string{"id", "name", "path", "version", "pinned_version", "issue"} {
		requireField(t, status, field, "string")
	}
	requireField(t, status, "valid", "bool")
}

func TestLockResultSchemaContract(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	if err := WriteLock(&buffer, testEnvelopeMetadata("lock"), FormatJSON, LockResult{
		Path:         "/repo/quill.lock",
		ArchiveCount: 3,
	}); err != nil {
		t.Fatalf("WriteLock: %v", err)
	}
	result := decodeResult(t, buffer.Bytes())

	requireField(t, result, "path", "string")
	requireField(t, result, "archive_count", "number")
}

func TestListResultSchemaContract(t *testing.T) {
	t.Parallel()

	t.Run("packs", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		if err := WriteList(&buffer, testEnvelopeMetadata("list"), FormatJSON, ListResult{
			Selector: ListPacks,
			Packs: []ListPack{{
				ID: "project", Name: "Project", Active: true, Provenance: "shipped", Rules: 2, Tools: 1,
			}},
		}); err != nil {
			t.Fatalf("WriteList: %v", err)
		}
		pack := requireArray(t, decodeResult(t, buffer.Bytes()), "packs")[0].(map[string]any)
		for _, field := range []string{"id", "name", "provenance"} {
			requireField(t, pack, field, "string")
		}
		requireField(t, pack, "active", "bool")
		requireField(t, pack, "rules", "number")
		requireField(t, pack, "tools", "number")
	})

	t.Run("rules", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		if err := WriteList(&buffer, testEnvelopeMetadata("list"), FormatJSON, ListResult{
			Selector: ListRules,
			Rules: []ListRule{{
				ID: "gofmt", Pack: "project", Provenance: "shipped", Name: "gofmt",
				Active: true, Enforcement: "required", Scope: "all", Fix: true,
			}},
		}); err != nil {
			t.Fatalf("WriteList: %v", err)
		}
		rule := requireArray(t, decodeResult(t, buffer.Bytes()), "rules")[0].(map[string]any)
		for _, field := range []string{"id", "pack", "provenance", "name", "enforcement", "scope"} {
			requireField(t, rule, field, "string")
		}
		requireField(t, rule, "active", "bool")
		requireField(t, rule, "fix", "bool")
	})

	t.Run("tools", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		if err := WriteList(&buffer, testEnvelopeMetadata("list"), FormatJSON, ListResult{
			Selector: ListTools,
			Tools: []ListTool{{
				ID: "gofmt", Name: "gofmt", Command: "gofmt", Pin: "1.24.5", Packs: []string{"project"},
			}},
		}); err != nil {
			t.Fatalf("WriteList: %v", err)
		}
		tool := requireArray(t, decodeResult(t, buffer.Bytes()), "tools")[0].(map[string]any)
		for _, field := range []string{"id", "name", "command", "pin"} {
			requireField(t, tool, field, "string")
		}
		requireField(t, tool, "packs", "array")
	})

	t.Run("scopes", func(t *testing.T) {
		t.Parallel()

		var buffer bytes.Buffer
		if err := WriteList(&buffer, testEnvelopeMetadata("list"), FormatJSON, ListResult{
			Selector: ListScopes,
			Scopes: []ListScope{{
				Name: "all", Roots: []string{"."}, Default: true,
			}},
		}); err != nil {
			t.Fatalf("WriteList: %v", err)
		}
		scope := requireArray(t, decodeResult(t, buffer.Bytes()), "scopes")[0].(map[string]any)
		requireField(t, scope, "name", "string")
		requireField(t, scope, "roots", "array")
		requireField(t, scope, "default", "bool")
	})
}

func TestExplainResultSchemaContract(t *testing.T) {
	t.Parallel()

	fix := ExplainExecution{
		Category: "file_command", Tools: []string{"gofmt"}, FileSet: "*.go", Language: "go",
	}
	result := ExplainResult{Rule: ExplainRule{
		ID:    "gofmt",
		Name:  "gofmt",
		Group: "formatting",
		Pack: ExplainPack{
			ID: "project", Name: "Project", Provenance: "shipped",
			Policy: map[string]any{"line-limit": float64(100)},
		},
		Binding: ExplainBinding{Enforcement: "required", Scope: "all", Requirements: []string{"1.1.format"}},
		Check:   ExplainExecution{Category: "file_command", Tools: []string{"gofmt"}},
		Fix:     &fix,
	}}

	var buffer bytes.Buffer
	if err := WriteExplain(&buffer, testEnvelopeMetadata("explain"), FormatJSON, result); err != nil {
		t.Fatalf("WriteExplain: %v", err)
	}
	rule := requireObject(t, decodeResult(t, buffer.Bytes()), "rule")

	for _, field := range []string{"id", "name", "group"} {
		requireField(t, rule, field, "string")
	}

	pack := requireObject(t, rule, "pack")
	for _, field := range []string{"id", "name", "provenance"} {
		requireField(t, pack, field, "string")
	}
	requireField(t, pack, "policy", "object")

	binding := requireObject(t, rule, "binding")
	for _, field := range []string{"enforcement", "scope"} {
		requireField(t, binding, field, "string")
	}
	requireField(t, binding, "requirements", "array")

	check := requireObject(t, rule, "check")
	requireField(t, check, "category", "string")
	requireField(t, rule, "fix", "object")
}

// TestCheckPerRuleExecutionErrorIsEnvelopeOK defends the protocol classification documented in
// docs/cli-protocol.md: a rule whose check ran but could not complete is a per-entry finding
// (status "error" with execution_error), never a command-level operation_failed envelope. The
// command completed, so the envelope stays status ok; the documented exit status 1 is the CLI's
// separate responsibility.
func TestCheckPerRuleExecutionErrorIsEnvelopeOK(t *testing.T) {
	t.Parallel()

	view := NewCheckView(CheckResult{
		Entries: []CheckEntry{{
			Rule: NewRuleSummary(style.Rule{
				ID:          "gofmt",
				Name:        "gofmt",
				Group:       style.RuleGroup("formatting"),
				Enforcement: style.EnforcementRequired,
				Scope:       style.Scope("all"),
			}),
			Status:         style.CheckStatusError,
			ExecutionError: errors.New("gofmt not installed"),
		}},
	})

	var buffer bytes.Buffer
	if _, err := WriteCheck(&buffer, testEnvelopeMetadata("check"), FormatJSON, view, false); err != nil {
		t.Fatalf("WriteCheck: %v", err)
	}

	var envelope struct {
		Status string `json:"status"`
		Result struct {
			Result struct {
				Entries []struct {
					Status         string `json:"status"`
					ExecutionError string `json:"execution_error"`
				} `json:"entries"`
			} `json:"result"`
			Summary CheckSummary `json:"summary"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode check envelope: %v", err)
	}

	if envelope.Status != StatusOK {
		t.Fatalf("per-rule execution error must keep envelope status %q, got %q", StatusOK, envelope.Status)
	}
	if envelope.Result.Summary.Errored != 1 {
		t.Fatalf("summary Errored = %d, want 1", envelope.Result.Summary.Errored)
	}
	if len(envelope.Result.Result.Entries) != 1 {
		t.Fatalf("expected one entry, got %d", len(envelope.Result.Result.Entries))
	}
	entry := envelope.Result.Result.Entries[0]
	if entry.Status != string(style.CheckStatusError) {
		t.Fatalf("entry status = %q, want %q", entry.Status, style.CheckStatusError)
	}
	if entry.ExecutionError == "" {
		t.Fatal("an errored entry must carry a non-empty execution_error")
	}
}
