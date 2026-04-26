package report

import (
	"fmt"
	"io"
	"strings"

	"ciphera/tools/internal/contract"
)

/* ---------------------------------------- Check Output ---------------------------------------- */

func WriteCheck(
	writer io.Writer,
	format OutputFormat,
	view CheckView,
	verbose bool,
) (summary CheckSummary, err error) {
	switch format {
	case FormatText:
		return writeCheckText(writer, view, verbose)
	case FormatJSON:
		return writeCheckJSON(writer, view)
	default:
		return summary, fmt.Errorf("unsupported output format %q", format)
	}
}

func writeCheckText(
	writer io.Writer,
	view CheckView,
	verbose bool,
) (summary CheckSummary, err error) {
	summary = view.Summary

	if _, err = fmt.Fprintln(writer, ""); err != nil {
		return summary, err
	}

	if _, err = fmt.Fprintln(writer, "STYLE.md Compliance Check"); err != nil {
		return summary, err
	}

	if _, err = fmt.Fprintln(writer, ""); err != nil {
		return summary, err
	}

	for groupIndex, group := range view.Groups {
		if groupIndex > 0 {
			if _, err = fmt.Fprintln(writer, ""); err != nil {
				return summary, err
			}
		}

		if _, err = fmt.Fprintln(writer, groupLabel(group.Group)); err != nil {
			return summary, err
		}

		for _, entry := range group.Entries {
			if err = writeAlignedColumns(
				writer,
				"  ["+entry.Rule.ID+"]",
				entry.Rule.Name,
				strings.ToUpper(string(entry.Status)),
			); err != nil {
				return summary, err
			}

			if err = writeRuleDetails(writer, entry, verbose); err != nil {
				return summary, err
			}
		}
	}

	if _, err = fmt.Fprintln(writer, ""); err != nil {
		return summary, err
	}

	_, err = fmt.Fprintf(
		writer,
		"Results: %d passed, %d warned, %d failed, %d skipped\n",
		summary.Passed,
		summary.Warned,
		summary.Failed,
		summary.Skipped,
	)
	return summary, err
}

func writeCheckJSON(writer io.Writer, view CheckView) (summary CheckSummary, err error) {
	summary = view.Summary
	err = writeJSON(writer, struct {
		Check checkJSON `json:"check"`
	}{Check: newCheckJSON(view)})
	return summary, err
}

func writeRuleDetails(writer io.Writer, entry CheckEntry, verbose bool) (err error) {
	if entry.Status == contract.CheckStatusPass {
		return nil
	}

	if verbose && len(entry.Rule.RequirementIDs) > 0 {
		if _, err = fmt.Fprintf(
			writer,
			"    requirements: %s\n",
			strings.Join(entry.Rule.RequirementIDs, ", "),
		); err != nil {
			return err
		}
	}

	if !verbose {
		return nil
	}

	for _, diagnostic := range entry.Result.Diagnostics {
		if _, err = fmt.Fprintf(writer, "    %s\n", formatDiagnostic(diagnostic)); err != nil {
			return err
		}
	}

	if strings.TrimSpace(entry.Result.Output) != "" {
		for _, line := range strings.Split(strings.TrimSpace(entry.Result.Output), "\n") {
			if _, err = fmt.Fprintf(writer, "    %s\n", line); err != nil {
				return err
			}
		}
	}

	if commandMetadataPresent(entry.Result.Command) {
		if _, err = fmt.Fprintf(
			writer,
			"    command: exit_code=%d timed_out=%t truncated=%t\n",
			entry.Result.Command.ExitCode,
			entry.Result.Command.TimedOut,
			entry.Result.Command.Truncated,
		); err != nil {
			return err
		}
	}

	return nil
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
			Level:        entry.Rule.Level,
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

func diagnosticListJSON(diagnostics []contract.Diagnostic) (payload []diagnosticJSON) {
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

func commandResultJSONFor(command contract.CommandResult) (payload *commandResultJSON) {
	if !commandMetadataPresent(command) {
		return nil
	}

	return &commandResultJSON{
		ExitCode:  command.ExitCode,
		TimedOut:  command.TimedOut,
		Truncated: command.Truncated,
	}
}

func commandMetadataPresent(command contract.CommandResult) (present bool) {
	return command != contract.CommandResult{}
}

func formatDiagnostic(diagnostic contract.Diagnostic) (line string) {
	location := diagnostic.File
	if diagnostic.Line > 0 {
		location = fmt.Sprintf("%s:%d", location, diagnostic.Line)
		if diagnostic.Column > 0 {
			location = fmt.Sprintf("%s:%d", location, diagnostic.Column)
		}
	}

	if diagnostic.Code == "" {
		return fmt.Sprintf("%s %s", location, diagnostic.Message)
	}

	return fmt.Sprintf("%s: [%s] %s", location, diagnostic.Code, diagnostic.Message)
}

func groupLabel(group contract.RuleGroup) (label string) {
	words := strings.FieldsFunc(string(group), func(character rune) bool {
		return character == '_' || character == '-' || character == '/'
	})
	for index, word := range words {
		if word == "" {
			continue
		}

		words[index] = strings.ToUpper(word[:1]) + word[1:]
	}

	return strings.Join(words, " ")
}
