package report

import (
	"fmt"
	"io"
)

// WriteCoverage writes coverage in the requested format. In JSON mode it writes the full
// machine envelope identified by metadata.
func WriteCoverage(
	writer io.Writer,
	metadata EnvelopeMetadata,
	format OutputFormat,
	view CoverageView,
	verbose bool,
) (err error) {
	switch format {
	case FormatText:
		return writeCoverageText(writer, view, verbose)
	case FormatJSON:
		return writeCoverageJSON(writer, metadata, view)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}
