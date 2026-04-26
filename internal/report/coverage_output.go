package report

import (
	"fmt"
	"io"
	"strings"

	"ciphera/tools/internal/coverage"
)

/* --------------------------------------- Coverage Output -------------------------------------- */

func WriteCoverage(
	writer io.Writer,
	format OutputFormat,
	view CoverageView,
	verbose bool,
) (err error) {
	switch format {
	case FormatText:
		return writeCoverageText(writer, view, verbose)
	case FormatJSON:
		return writeCoverageJSON(writer, view)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

func writeCoverageText(writer io.Writer, view CoverageView, verbose bool) (err error) {
	lines := []string{
		"STYLE.md Coverage",
		"",
		"Requirements",
		fmt.Sprintf("Automated:          %d", view.RequirementTotals.Automated),
		fmt.Sprintf("Review-only:        %d", view.RequirementTotals.ReviewOnly),
		fmt.Sprintf("Manual-deferred:    %d", view.RequirementTotals.ManualDeferred),
		"",
		"Sections",
		fmt.Sprintf("Automated:   %d", view.SectionTotals.Automated),
		fmt.Sprintf("Partial:     %d", view.SectionTotals.Partial),
		fmt.Sprintf("Review-only: %d", view.SectionTotals.ReviewOnly),
		fmt.Sprintf("Manual:      %d", view.SectionTotals.Manual),
		"",
	}
	for _, line := range lines {
		if _, err = fmt.Fprintln(writer, line); err != nil {
			return err
		}
	}

	for _, entry := range view.Report.Sections {
		if err = writeAlignedColumns(
			writer,
			"["+entry.Section+"]",
			entry.Title,
			strings.ToUpper(string(entry.Status)),
			coverageSummary(entry),
		); err != nil {
			return err
		}
	}

	if !verbose {
		return nil
	}

	return writeCoverageDetails(writer, view.Outstanding)
}

func writeCoverageDetails(writer io.Writer, requirements []coverage.Requirement) (err error) {
	if len(requirements) == 0 {
		return nil
	}

	if _, err = fmt.Fprintln(writer, ""); err != nil {
		return err
	}

	if _, err = fmt.Fprintln(writer, "Outstanding Requirements"); err != nil {
		return err
	}

	for _, requirement := range requirements {
		if err = writeAlignedColumns(
			writer,
			"["+requirement.Section+"]",
			strings.ToUpper(string(requirement.Mode)),
			requirement.ID,
		); err != nil {
			return err
		}

		if _, err = fmt.Fprintf(writer, "    %s\n", requirement.Text); err != nil {
			return err
		}

		if requirement.Reason == "" {
			continue
		}

		if _, err = fmt.Fprintf(writer, "    why: %s\n", requirement.Reason); err != nil {
			return err
		}
	}

	return nil
}

func coverageSummary(entry coverage.Section) (summary string) {
	parts := []string{fmt.Sprintf("%d/%d automated", entry.AutomatedCount, entry.RequirementCount)}
	if entry.ManualDeferredCount > 0 {
		parts = append(parts, fmt.Sprintf("%d deferred", entry.ManualDeferredCount))
	}

	return "(" + strings.Join(parts, ", ") + ")"
}

func writeCoverageJSON(writer io.Writer, view CoverageView) (err error) {
	return writeJSON(writer, struct {
		Coverage coverageJSON `json:"coverage"`
	}{Coverage: newCoverageJSON(view)})
}

func newCoverageJSON(view CoverageView) (payload coverageJSON) {
	return coverageJSON{
		Report: coverageReportJSON{
			Requirements: requirementListJSON(view.Report.Requirements),
			Sections:     sectionListJSON(view.Report.Sections),
		},
		RequirementTotals: view.RequirementTotals,
		SectionTotals:     view.SectionTotals,
		Outstanding:       requirementListJSON(view.Outstanding),
		OutstandingByMode: cloneIntMap(view.OutstandingByMode),
	}
}

func requirementListJSON(requirements []coverage.Requirement) (payload []requirementJSON) {
	payload = make([]requirementJSON, 0, len(requirements))
	for _, requirement := range requirements {
		payload = append(payload, requirementJSON{
			ID:      requirement.ID,
			Section: requirement.Section,
			Text:    requirement.Text,
			Mode:    string(requirement.Mode),
			Reason:  requirement.Reason,
			RuleIDs: append([]string{}, requirement.RuleIDs...),
		})
	}

	return payload
}

func sectionListJSON(sections []coverage.Section) (payload []sectionJSON) {
	payload = make([]sectionJSON, 0, len(sections))
	for _, section := range sections {
		payload = append(payload, sectionJSON{
			Section:             section.Section,
			Title:               section.Title,
			Status:              section.Status,
			RequirementCount:    section.RequirementCount,
			AutomatedCount:      section.AutomatedCount,
			ReviewOnlyCount:     section.ReviewOnlyCount,
			ManualDeferredCount: section.ManualDeferredCount,
		})
	}

	return payload
}

func cloneIntMap(source map[string]int) (target map[string]int) {
	target = make(map[string]int, len(source))
	for key, value := range source {
		target[key] = value
	}

	return target
}
