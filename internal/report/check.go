package report

import (
	"fmt"
	"io"
)

// WriteCheck write check.
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
