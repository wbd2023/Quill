package report

import (
	"io"

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

// WriteFix writes the machine-mode fix result envelope for command. Text-mode fix output is
// owned by the CLI, so this renderer only emits the JSON envelope.
func WriteFix(writer io.Writer, command string, view FixView) (err error) {
	return writeResultEnvelope(writer, command, newFixJSON(view))
}
