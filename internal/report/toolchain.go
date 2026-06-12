package report

import (
	"fmt"
	"io"

	"ciphera/tools/internal/toolchain"
)

type ToolchainResult struct {
	Statuses []toolchain.Status
}

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
