package report

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
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
			pack.Provenance,
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

	header := []string{
		"ID", "Pack", "Provenance", "Name", "State", "Enforcement", "Scope", "Fix",
	}
	rows := make([][]string, 0, len(rules))
	for _, rule := range rules {
		rows = append(rows, []string{
			rule.ID,
			rule.Pack,
			rule.Provenance,
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

	header := []string{"ID", "Name", "Pin", "Packs"}
	rows := make([][]string, 0, len(tools))
	for _, tool := range tools {
		rows = append(rows, []string{
			tool.ID,
			tool.Name,
			tool.Pin,
			strings.Join(tool.Packs, ","),
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

func yesNo(value bool) (label string) {
	if value {
		return "yes"
	}
	return "no"
}
