package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const tabwriterPadding = 2

func writeJSON(writer io.Writer, value any) (err error) {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func writeAlignedColumns(writer io.Writer, columns ...string) (err error) {
	var buffer bytes.Buffer

	table := tabwriter.NewWriter(&buffer, 0, 0, tabwriterPadding, ' ', 0)
	if _, err = fmt.Fprintln(table, strings.Join(columns, "\t")); err != nil {
		return err
	}

	if err = table.Flush(); err != nil {
		return err
	}

	_, err = buffer.WriteTo(writer)
	return err
}
