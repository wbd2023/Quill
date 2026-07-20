package report

import (
	"fmt"
	"io"

	"github.com/wbd2023/Quill/internal/toolchain"
)

// ToolchainResult is toolchain result.
type ToolchainResult struct {
	Statuses []toolchain.Status
}

// WriteToolchain write toolchain.
func WriteToolchain(
	writer io.Writer,
	format OutputFormat,
	view ToolchainView,
) (allValid bool, err error) {
	switch format {
	case FormatText:
		return writeToolchainText(writer, view)
	case FormatJSON:
		return writeToolchainJSON(writer, view)
	default:
		return false, fmt.Errorf("unsupported output format %q", format)
	}
}
