package report

import (
	"io"
	"strings"

	"ciphera/tools/internal/style"
)

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
			Output:       strings.TrimSpace(entry.Result.Output),
			Command:      commandResultJSONFor(entry.Result.Command),
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

func commandResultJSONFor(command style.CommandResult) (payload *commandResultJSON) {
	if !commandMetadataPresent(command) {
		return nil
	}

	return &commandResultJSON{
		ExitCode:  command.ExitCode,
		TimedOut:  command.TimedOut,
		Truncated: command.Truncated,
	}
}

func commandMetadataPresent(command style.CommandResult) (present bool) {
	return command != style.CommandResult{}
}
