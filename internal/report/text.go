package report

import (
	"fmt"
	"io"
	"text/tabwriter"
)

const (
	tablePadding = 2
	indentPrefix = "  "
)

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
