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
		Check CheckView `json:"check"`
	}{Check: view})
	return summary, err
}

func writeRuleDetails(writer io.Writer, entry CheckEntry, verbose bool) (err error) {
	if entry.Status == CheckStatusPass {
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

	if !verbose || strings.TrimSpace(entry.Output) == "" {
		return nil
	}

	for _, line := range strings.Split(strings.TrimSpace(entry.Output), "\n") {
		if _, err = fmt.Fprintf(writer, "    %s\n", line); err != nil {
			return err
		}
	}

	return nil
}

func groupLabel(group contract.RuleGroup) (label string) {
	switch group {
	case contract.RuleGroupControlPlane:
		return "Control Plane"

	case contract.RuleGroupLanguage:
		return "Language Backends"

	case contract.RuleGroupRepository:
		return "Repository Scanners"

	case contract.RuleGroupExternal:
		return "External Tools"

	default:
		return string(group)
	}
}
