package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

const testQuillVersion = "test-version"

func testEnvelopeMetadata(command string) (metadata EnvelopeMetadata) {
	return EnvelopeMetadata{Command: command, QuillVersion: testQuillVersion}
}

/* ---------------------------------------- Fix Envelope ---------------------------------------- */

func TestWriteFixDecodedEnvelope(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	view := NewFixView(
		style.Scope("tools"),
		false,
		[]toolchain.Status{
			{
				Tool:    toolchain.Tool{ID: "markdownlint", Name: "markdownlint"},
				Version: "0.48.0",
				Valid:   false,
				Issue:   "requires pinned version 0.45.0",
			},
		},
		[]FixEntry{
			{
				Rule: NewRuleSummary(style.Rule{
					ID:          "markdown",
					Name:        "markdownlint",
					Group:       style.RuleGroup("external_tools"),
					Enforcement: style.EnforcementRequired,
					Scope:       style.Scope("tools"),
				}),
				Execution: style.ExecutionResult{
					ExitCode:  2,
					TimedOut:  false,
					Truncated: true,
				},
				ExecutionError: errors.New("formatter exited non-zero"),
			},
		},
	)

	if _, err := WriteFix(&buffer, testEnvelopeMetadata("fix"), FormatJSON, view); err != nil {
		t.Fatalf("WriteFix: %v", err)
	}

	var envelope struct {
		SchemaVersion int    `json:"schema_version"`
		Command       string `json:"command"`
		Status        string `json:"status"`
		Result        struct {
			Scope     string `json:"scope"`
			Toolchain struct {
				AllValid bool `json:"all_valid"`
				Statuses []struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Valid bool   `json:"valid"`
				} `json:"statuses"`
			} `json:"toolchain"`
			Rules []struct {
				RuleID         string `json:"rule_id"`
				ExitCode       int    `json:"exit_code"`
				Truncated      bool   `json:"truncated"`
				ExecutionError string `json:"execution_error"`
			} `json:"rules"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode fix json: %v", err)
	}

	if envelope.SchemaVersion != SchemaVersion || envelope.Command != "fix" ||
		envelope.Status != StatusOK {
		t.Fatalf("unexpected envelope header: %+v", envelope)
	}

	if envelope.Result.Scope != "tools" {
		t.Fatalf("unexpected scope: %q", envelope.Result.Scope)
	}

	if envelope.Result.Toolchain.AllValid {
		t.Fatal("expected all_valid=false")
	}

	if len(envelope.Result.Toolchain.Statuses) != 1 ||
		envelope.Result.Toolchain.Statuses[0].ID != "markdownlint" {
		t.Fatalf("unexpected toolchain statuses: %+v",
			envelope.Result.Toolchain.Statuses)
	}

	if len(envelope.Result.Rules) != 1 {
		t.Fatalf("expected one rule, got %d", len(envelope.Result.Rules))
	}

	rule := envelope.Result.Rules[0]
	if rule.RuleID != "markdown" || rule.ExitCode != 2 || !rule.Truncated ||
		rule.ExecutionError != "formatter exited non-zero" {
		t.Fatalf("unexpected rule payload: %+v", rule)
	}
}

func TestFixViewHasExecutionError(t *testing.T) {
	t.Parallel()

	clean := NewFixView(style.Scope("tools"), true, nil, []FixEntry{
		{Rule: NewRuleSummary(style.Rule{ID: "clean"})},
	})
	if clean.HasExecutionError() {
		t.Fatal("clean view reported an execution error")
	}

	failing := NewFixView(style.Scope("tools"), true, nil, []FixEntry{
		{Rule: NewRuleSummary(style.Rule{ID: "failing"}), ExecutionError: errors.New("boom")},
	})
	if !failing.HasExecutionError() {
		t.Fatal("failing view did not report an execution error")
	}
}

/* ---------------------------------------- Lock Envelope --------------------------------------- */

func TestWriteLockDecodedEnvelope(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	if err := WriteLock(&buffer, testEnvelopeMetadata("lock"), FormatJSON, LockResult{
		Path:         "/repo/quill.lock",
		ArchiveCount: 3,
	}); err != nil {
		t.Fatalf("WriteLock: %v", err)
	}

	var envelope struct {
		SchemaVersion int    `json:"schema_version"`
		Command       string `json:"command"`
		Status        string `json:"status"`
		Result        struct {
			Path         string `json:"path"`
			ArchiveCount int    `json:"archive_count"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode lock json: %v", err)
	}

	if envelope.SchemaVersion != SchemaVersion || envelope.Command != "lock" ||
		envelope.Status != StatusOK {
		t.Fatalf("unexpected envelope header: %+v", envelope)
	}

	if envelope.Result.Path != "/repo/quill.lock" || envelope.Result.ArchiveCount != 3 {
		t.Fatalf("unexpected lock payload: %+v", envelope.Result)
	}
}

/* --------------------------------------- Error Envelope --------------------------------------- */

func TestNewErrorEnvelopeShape(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	if err := WriteEnvelope(&buffer, NewErrorEnvelope(
		testEnvelopeMetadata("check"),
		ErrorCodeInvalidArgument,
		"--scope must name a configured scope",
	)); err != nil {
		t.Fatalf("WriteEnvelope: %v", err)
	}

	var envelope Envelope
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode error envelope: %v", err)
	}

	if envelope.SchemaVersion != SchemaVersion {
		t.Fatalf("unexpected schema version: %d", envelope.SchemaVersion)
	}

	if envelope.QuillVersion != testQuillVersion {
		t.Fatalf("unexpected Quill version: %q", envelope.QuillVersion)
	}

	if envelope.Command != "check" || envelope.Status != StatusError {
		t.Fatalf("unexpected error header: command=%q status=%q",
			envelope.Command, envelope.Status)
	}

	if envelope.Result != nil {
		t.Fatalf("error envelope must not carry a result, got %v", envelope.Result)
	}

	if envelope.Error == nil || envelope.Error.Code != ErrorCodeInvalidArgument ||
		envelope.Error.Message != "--scope must name a configured scope" {
		t.Fatalf("unexpected error payload: %+v", envelope.Error)
	}
}

func TestNewResultEnvelopeHasNoError(t *testing.T) {
	t.Parallel()

	envelope := NewResultEnvelope(testEnvelopeMetadata("coverage"), map[string]int{"automated": 1})
	if envelope.Status != StatusOK || envelope.Error != nil {
		t.Fatalf("result envelope must be ok with no error: %+v", envelope)
	}

	if envelope.QuillVersion != testQuillVersion {
		t.Fatalf("unexpected Quill version: %q", envelope.QuillVersion)
	}
}
