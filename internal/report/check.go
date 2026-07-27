package report

import (
	"fmt"
	"io"
)

// WriteCheck writes a check result in the requested format. In JSON mode it writes the full
// machine envelope tagged with command.
func WriteCheck(
	writer io.Writer,
	command string,
	format OutputFormat,
	view CheckView,
	verbose bool,
) (summary CheckSummary, err error) {
	switch format {
	case FormatText:
		return writeCheckText(writer, view, verbose)
	case FormatJSON:
		return writeCheckJSON(writer, command, view)
	default:
		return summary, fmt.Errorf("unsupported output format %q", format)
	}
}
