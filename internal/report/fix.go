package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

/* ----------------------------------------- Fix Result ----------------------------------------- */

// FixEntry captures the outcome of one attempted fixer.
type FixEntry struct {
	Rule           RuleSummary
	Execution      style.ExecutionResult
	ExecutionError error
}

// FixView is the presentation view for a completed fix operation.
type FixView struct {
	Scope    style.Scope
	AllValid bool
	Statuses []toolchain.Status
	Entries  []FixEntry
}

// NewFixView builds a fix view from the resolved scope, toolchain inspection, and attempted
// fixers.
func NewFixView(
	scope style.Scope,
	allValid bool,
	statuses []toolchain.Status,
	entries []FixEntry,
) (view FixView) {
	return FixView{
		Scope:    scope,
		AllValid: allValid,
		Statuses: statuses,
		Entries:  entries,
	}
}

// NewFixResult converts a completed engine fix into the presentation view used by both text and
// JSON renderers.
func NewFixResult(result engine.FixResult) (view FixView) {
	entries := make([]FixEntry, 0, len(result.Rules))
	for _, rule := range result.Rules {
		entries = append(entries, FixEntry{
			Rule:           NewRuleSummary(rule.Rule),
			Execution:      rule.Execution,
			ExecutionError: rule.ExecutionError,
		})
	}

	return NewFixView(
		result.Scope,
		result.Toolchain.AllValid,
		result.Toolchain.Statuses,
		entries,
	)
}

// HasExecutionError reports whether any fixer failed during execution. The engine records
// per-fixer failures in the result rather than aborting the operation, so the envelope reports
// status "ok" with the failures embedded and the CLI maps them to a nonzero exit status.
func (view FixView) HasExecutionError() (hasError bool) {
	for _, entry := range view.Entries {
		if entry.ExecutionError != nil {
			return true
		}
	}

	return false
}

// FixSummary contains the status policy inputs for a completed fix render.
type FixSummary struct {
	AllValid          bool
	HasExecutionError bool
}

// WriteFix writes a fix result in the requested format. Text diagnostic output is intentionally
// routed by the CLI to stderr, while JSON output is the sole stdout document.
func WriteFix(
	writer io.Writer,
	metadata EnvelopeMetadata,
	format OutputFormat,
	view FixView,
) (summary FixSummary, err error) {
	summary.AllValid = view.AllValid
	summary.HasExecutionError = view.HasExecutionError()

	switch format {
	case FormatText:
		return writeFixText(writer, view, summary)
	case FormatJSON:
		return summary, writeResultEnvelope(writer, metadata, newFixJSON(view))
	default:
		return FixSummary{}, fmt.Errorf("unsupported output format %q", format)
	}
}

func writeFixText(
	writer io.Writer,
	view FixView,
	summary FixSummary,
) (result FixSummary, err error) {
	if len(view.Entries) == 0 {
		return summary, nil
	}

	if _, err := writeToolchainText(writer, ToolchainResult{
		AllValid: view.AllValid,
		Statuses: view.Statuses,
	}); err != nil {
		return FixSummary{}, err
	}

	if !summary.AllValid {
		return summary, nil
	}

	for _, entry := range view.Entries {
		if entry.ExecutionError == nil {
			continue
		}
		if output := strings.TrimSpace(entry.Execution.Output); output != "" {
			_, err := fmt.Fprintln(writer, output)
			return summary, err
		}
		_, err := fmt.Fprintln(writer, entry.ExecutionError)
		return summary, err
	}

	return summary, nil
}
