package report

import (
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/toolchain"
)

// ToolchainResult is toolchain result.
type ToolchainResult struct {
	Statuses []toolchain.Status
}

// WriteToolchain writes a toolchain inspection in the requested format. In JSON mode it writes
// the full machine envelope tagged with command (used by both doctor and install).
func WriteToolchain(
	writer io.Writer,
	command string,
	format OutputFormat,
	view ToolchainView,
) (allValid bool, err error) {
	switch format {
	case FormatText:
		return writeToolchainText(writer, view)
	case FormatJSON:
		return writeToolchainJSON(writer, command, view)
	default:
		return false, fmt.Errorf("unsupported output format %q", format)
	}
}
