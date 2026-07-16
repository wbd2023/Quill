package report

import (
	"io"

	"ciphera/tools/internal/style"
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
	RuleID       string             `json:"rule_id"`
	Name         string             `json:"name"`
	Group        style.RuleGroup    `json:"group"`
	Enforcement  style.Enforcement  `json:"enforcement"`
	Scope        style.Scope        `json:"scope"`
	Status       style.CheckStatus  `json:"status"`
	Requirements []string           `json:"requirements"`
	Diagnostics  []diagnosticJSON   `json:"diagnostics"`
	Command      *commandResultJSON `json:"command,omitempty"`
}

type diagnosticJSON struct {
	Code    string `json:"code"`
	File    string `json:"file,omitempty"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	Message string `json:"message"`
}

type commandResultJSON struct {
	ExitCode  int  `json:"exit_code"`
	TimedOut  bool `json:"timed_out"`
	Truncated bool `json:"truncated"`
}

/* ------------------------------------------ Rendering ----------------------------------------- */

func writeCheckJSON(writer io.Writer, view CheckView) (summary CheckSummary, err error) {
	summary = view.Summary
	err = writeJSON(writer, struct {
		Check checkJSON `json:"check"`
	}{Check: newCheckJSON(view)})
	return summary, err
}

func newCheckJSON(view CheckView) (payload checkJSON) {
	payload = checkJSON{
		Result: checkResultJSON{
			Entries: checkEntryListJSON(view.Result.Entries),
		},
		Summary: view.Summary,
		Groups:  make([]checkGroupJSON, 0, len(view.Groups)),
	}

	for _, group := range view.Groups {
		payload.Groups = append(payload.Groups, checkGroupJSON{
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
			RuleID:       entry.Rule.ID,
			Name:         entry.Rule.Name,
			Group:        entry.Rule.Group,
			Enforcement:  entry.Rule.Enforcement,
			Scope:        entry.Rule.Scope,
			Status:       entry.Status,
			Requirements: append([]string{}, entry.Rule.RequirementIDs...),
			Diagnostics:  diagnosticListJSON(entry.Result.Diagnostics),
			Command:      commandResultJSONFor(entry.Result),
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
			Line:    diagnostic.Line,
			Column:  diagnostic.Column,
			Message: diagnostic.Message,
		})
	}

	return payload
}

func commandResultJSONFor(result style.ExecutionResult) (payload *commandResultJSON) {
	if !result.HasCommand() {
		return nil
	}

	return &commandResultJSON{
		ExitCode:  result.ExitCode,
		TimedOut:  result.TimedOut,
		Truncated: result.Truncated,
	}
}
