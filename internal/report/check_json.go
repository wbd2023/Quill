package report

import (
	"io"

	"github.com/wbd2023/quill/internal/style"
)

/* ------------------------------------------ JSON DTOs ----------------------------------------- */

type checkJSON struct {
	Result  checkResultJSON  `json:"result"`
	Summary CheckSummary     `json:"summary"`
	Groups  []checkGroupJSON `json:"groups"`
}

type checkResultJSON struct {
	Entries []checkEntryJSON `json:"entries"`
}

type checkGroupJSON struct {
	Group   style.RuleGroup  `json:"group"`
	Entries []checkEntryJSON `json:"entries"`
}

type checkEntryJSON struct {
	RuleID         string             `json:"rule_id"`
	Name           string             `json:"name"`
	Group          style.RuleGroup    `json:"group"`
	Enforcement    style.Enforcement  `json:"enforcement"`
	Scope          style.Scope        `json:"scope"`
	Status         style.CheckStatus  `json:"status"`
	Requirements   []string           `json:"requirements"`
	Diagnostics    []diagnosticJSON   `json:"diagnostics"`
	ExecutionError string             `json:"execution_error,omitempty"`
	Command        *commandResultJSON `json:"command,omitempty"`
}

type diagnosticJSON struct {
	Code    string     `json:"code"`
	File    string     `json:"file,omitempty"`
	Range   *rangeJSON `json:"range,omitempty"`
	Message string     `json:"message"`
	HelpURL string     `json:"help_url,omitempty"`
}

type rangeJSON struct {
	Start positionJSON  `json:"start"`
	End   *positionJSON `json:"end,omitempty"`
}

type positionJSON struct {
	Line   int `json:"line"`
	Column int `json:"column,omitempty"`
}

type commandResultJSON struct {
	ExitCode  int  `json:"exit_code"`
	TimedOut  bool `json:"timed_out"`
	Truncated bool `json:"truncated"`
}

/* ------------------------------------------ Rendering ----------------------------------------- */

func writeCheckJSON(
	writer io.Writer,
	command string,
	view CheckView,
) (summary CheckSummary, err error) {
	err = writeResultEnvelope(writer, command, newCheckJSON(view))
	return view.Summary, err
}

func newCheckJSON(view CheckView) (payload checkJSON) {
	return checkJSON{
		Result:  checkResultJSON{Entries: checkEntryListJSON(view.Result.Entries)},
		Summary: view.Summary,
		Groups:  checkGroupListJSON(view.Groups),
	}
}

func checkGroupListJSON(groups []CheckGroup) (payload []checkGroupJSON) {
	payload = make([]checkGroupJSON, 0, len(groups))
	for _, group := range groups {
		payload = append(payload, checkGroupJSON{
			Group:   group.Group,
			Entries: checkEntryListJSON(group.Entries),
		})
	}

	return payload
}

func checkEntryListJSON(entries []CheckEntry) (payload []checkEntryJSON) {
	payload = make([]checkEntryJSON, 0, len(entries))
	for _, entry := range entries {
		payload = append(payload, checkEntryJSON{
			RuleID:         entry.Rule.ID,
			Name:           entry.Rule.Name,
			Group:          entry.Rule.Group,
			Enforcement:    entry.Rule.Enforcement,
			Scope:          entry.Rule.Scope,
			Status:         entry.Status,
			Requirements:   append([]string{}, entry.Rule.RequirementIDs...),
			Diagnostics:    diagnosticListJSON(entry.Result.Diagnostics),
			ExecutionError: executionErrorText(entry.ExecutionError),
			Command:        commandResultJSONFor(entry.Result),
		})
	}

	return payload
}

func diagnosticListJSON(diagnostics []style.Diagnostic) (payload []diagnosticJSON) {
	payload = make([]diagnosticJSON, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		payload = append(payload, diagnosticJSON{
			Code:    diagnostic.Code,
			File:    diagnostic.File,
			Range:   rangeJSONFor(diagnostic.Range),
			Message: diagnostic.Message,
			HelpURL: diagnostic.HelpURL,
		})
	}

	return payload
}

// rangeJSONFor converts an internal range to its persisted JSON form. The range is omitted when no
// start line is known, making absent locations explicit; End is included only when known.
func rangeJSONFor(location style.Range) (payload *rangeJSON) {
	if !location.Start.IsKnown() {
		return nil
	}

	payload = &rangeJSON{
		Start: positionJSON{Line: location.Start.Line, Column: location.Start.Column},
	}
	if location.End.IsKnown() {
		payload.End = &positionJSON{Line: location.End.Line, Column: location.End.Column}
	}

	return payload
}

func executionErrorText(err error) (text string) {
	if err == nil {
		return ""
	}

	return err.Error()
}

func commandResultJSONFor(result style.ExecutionResult) (payload *commandResultJSON) {
	if !hasCommandMetadata(result) {
		return nil
	}

	return &commandResultJSON{
		ExitCode:  result.ExitCode,
		TimedOut:  result.TimedOut,
		Truncated: result.Truncated,
	}
}
