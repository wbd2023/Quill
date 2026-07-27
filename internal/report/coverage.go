package report

import (
	"fmt"
	"io"
)

// WriteCoverage writes coverage in the requested format. In JSON mode it writes the full
// machine envelope tagged with command.
func WriteCoverage(
	writer io.Writer,
	command string,
	format OutputFormat,
	view CoverageView,
	verbose bool,
) (err error) {
	switch format {
	case FormatText:
		return writeCoverageText(writer, view, verbose)
	case FormatJSON:
		return writeCoverageJSON(writer, command, view)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}
