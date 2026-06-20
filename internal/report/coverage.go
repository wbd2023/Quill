package report

import (
	"fmt"
	"io"
)

// WriteCoverage write coverage.
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
