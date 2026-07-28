package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

/* ---------------------------------------- Text Helpers ---------------------------------------- */

const (
	tablePadding = 2
	indentPrefix = "  "
)

// writeTable writes a header row followed by data rows through one tabwriter so every column
// aligns across the whole table.
func writeTable(writer io.Writer, header []string, rows [][]string) (err error) {
	table := tabwriter.NewWriter(writer, 0, 0, tablePadding, ' ', 0)

	if err = writeTableRow(table, header); err != nil {
		return err
	}

	for _, row := range rows {
		if err = writeTableRow(table, row); err != nil {
			return err
		}
	}

	return table.Flush()
}

func writeTableRow(table *tabwriter.Writer, columns []string) (err error) {
	_, err = fmt.Fprintln(table, strings.Join(columns, "\t"))
	return err
}

/* ------------------------------------------ List Text ----------------------------------------- */

func writeListText(writer io.Writer, result ListResult) (err error) {
	switch result.Selector {
	case ListPacks:
		return writePacksText(writer, result.Packs)
	case ListRules:
		return writeRulesText(writer, result.Rules)
	case ListTools:
		return writeToolsText(writer, result.Tools)
	case ListScopes:
		return writeScopesText(writer, result.Scopes)
	default:
		return fmt.Errorf("unsupported list selector %q", result.Selector)
	}
}

func writePacksText(writer io.Writer, packs []ListPack) (err error) {
	if _, err = fmt.Fprintf(writer, "Packs (%d)\n\n", len(packs)); err != nil {
		return err
	}

	header := []string{"ID", "Name", "State", "Source", "Rules", "Tools"}
	rows := make([][]string, 0, len(packs))
	for _, pack := range packs {
		rows = append(rows, []string{
			pack.ID,
			pack.Name,
			packState(pack.Active),
			packSource(pack.External),
			fmt.Sprintf("%d", pack.Rules),
			fmt.Sprintf("%d", pack.Tools),
		})
	}

	return writeTable(writer, header, rows)
}

func writeRulesText(writer io.Writer, rules []ListRule) (err error) {
	if _, err = fmt.Fprintf(writer, "Rules (%d)\n\n", len(rules)); err != nil {
		return err
	}

	header := []string{"ID", "Pack", "Name", "State", "Enforcement", "Scope", "Fix"}
	rows := make([][]string, 0, len(rules))
	for _, rule := range rules {
		rows = append(rows, []string{
			rule.ID,
			rule.Pack,
			rule.Name,
			ruleState(rule.Active),
			rule.Enforcement,
			rule.Scope,
			yesNo(rule.Fix),
		})
	}

	return writeTable(writer, header, rows)
}

func writeToolsText(writer io.Writer, tools []ListTool) (err error) {
	if _, err = fmt.Fprintf(writer, "Tools (%d)\n\n", len(tools)); err != nil {
		return err
	}

	header := []string{"ID", "Name", "Pin", "Packs", "Source"}
	rows := make([][]string, 0, len(tools))
	for _, tool := range tools {
		rows = append(rows, []string{
			tool.ID,
			tool.Name,
			tool.Pin,
			strings.Join(tool.Packs, ","),
			packSource(tool.External),
		})
	}

	return writeTable(writer, header, rows)
}

func writeScopesText(writer io.Writer, scopes []ListScope) (err error) {
	if _, err = fmt.Fprintf(writer, "Scopes (%d)\n\n", len(scopes)); err != nil {
		return err
	}

	header := []string{"Name", "Default", "Roots"}
	rows := make([][]string, 0, len(scopes))
	for _, scope := range scopes {
		rows = append(rows, []string{
			scope.Name,
			yesNo(scope.Default),
			strings.Join(scope.Roots, ","),
		})
	}

	return writeTable(writer, header, rows)
}

/* ---------------------------------------- Explain Text ---------------------------------------- */

func writeExplainText(writer io.Writer, result ExplainResult) (err error) {
	rule := result.Rule

	if _, err = fmt.Fprintf(writer, "Rule: %s\n\n", rule.ID); err != nil {
		return err
	}

	rows := [][2]string{
		{"pack", fmt.Sprintf("%s (%s)", rule.Pack, packSource(rule.External))},
		{"name", rule.Name},
		{"group", rule.Group},
		{"enforcement", rule.Enforcement},
		{"scope", rule.Scope},
		{"requirements", strings.Join(rule.Requirements, ", ")},
	}
	rows = append(rows, executionRows("check", rule.Check)...)

	if rule.Fix == nil {
		rows = append(rows, [2]string{"fix", "none"})
	} else {
		rows = append(rows, executionRows("fix", *rule.Fix)...)
	}

	if err = writeFields(writer, rows); err != nil {
		return err
	}

	return writePackConfig(writer, rule.PackConfig)
}

// writePackConfig renders the relevant Pack configuration deterministically beneath the rule
// explanation. json.MarshalIndent emits map keys in sorted order, so the output is stable across
// runs and platforms. An empty or absent configuration renders nothing.
func writePackConfig(writer io.Writer, config map[string]any) (err error) {
	if len(config) == 0 {
		return nil
	}

	if _, err = fmt.Fprint(writer, "\nPack config\n\n"); err != nil {
		return err
	}

	encoded, err := json.MarshalIndent(config, "", indentPrefix)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(strings.TrimRight(string(encoded), "\n"), "\n") {
		if _, err = fmt.Fprintf(writer, "%s%s\n", indentPrefix, line); err != nil {
			return err
		}
	}

	return nil
}

func executionRows(label string, execution ExplainExecution) (rows [][2]string) {
	category := execution.Category
	if execution.Detail != "" {
		category = fmt.Sprintf("%s (%s)", execution.Category, execution.Detail)
	}

	rows = append(rows, [2]string{label, category})
	if len(execution.Tools) > 0 {
		rows = append(rows, [2]string{label + " tools", strings.Join(execution.Tools, ", ")})
	}
	if execution.FileSet != "" {
		rows = append(rows, [2]string{label + " file set", execution.FileSet})
	}
	if execution.Language != "" {
		rows = append(rows, [2]string{label + " language", execution.Language})
	}

	return rows
}

// writeFields renders label/value rows through one tabwriter so every label column aligns.
func writeFields(writer io.Writer, rows [][2]string) (err error) {
	table := tabwriter.NewWriter(writer, 0, 0, tablePadding, ' ', 0)
	for _, row := range rows {
		if _, err = fmt.Fprintf(table, "%s%s:\t%s\n", indentPrefix, row[0], row[1]); err != nil {
			return err
		}
	}

	return table.Flush()
}

/* ----------------------------------------- Formatting ----------------------------------------- */

func packState(active bool) (state string) {
	if active {
		return "active"
	}
	return "inactive"
}

func ruleState(active bool) (state string) {
	if active {
		return "active"
	}
	return "inactive"
}

func packSource(external bool) (source string) {
	if external {
		return "external"
	}
	return "built-in"
}

func yesNo(value bool) (label string) {
	if value {
		return "yes"
	}
	return "no"
}
