package report

import (
	"fmt"
	"io"
	"strings"

	"ciphera/tools/internal/style"
)

/* ----------------------------------------- Text Output ---------------------------------------- */

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
		"Results: %d passed, %d warned, %d failed, %d skipped, %d errored\n",
		summary.Passed,
		summary.Warned,
		summary.Failed,
		summary.Skipped,
		summary.Errored,
	)
	return summary, err
}

/* ---------------------------------------- Rule Details ---------------------------------------- */

func writeRuleDetails(writer io.Writer, entry CheckEntry, verbose bool) (err error) {
	if entry.Status == style.CheckStatusPass {
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

	if entry.Result.HasCommand() {
		if _, err = fmt.Fprintf(
			writer,
			"    command: exit_code=%d timed_out=%t truncated=%t\n",
			entry.Result.ExitCode,
			entry.Result.TimedOut,
			entry.Result.Truncated,
		); err != nil {
			return err
		}
	}

	return nil
}
