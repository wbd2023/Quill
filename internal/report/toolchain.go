package report

import (
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/toolchain"
)

// ToolchainResult is the presentation input for one engine toolchain inspection.
type ToolchainResult struct {
	AllValid bool
	Statuses []toolchain.Status
}

// NewToolchainResult converts one engine inspection into the shared toolchain presentation result.
func NewToolchainResult(inspection engine.ToolchainInspection) (result ToolchainResult) {
	return ToolchainResult{
		AllValid: inspection.AllValid,
		Statuses: inspection.Statuses,
	}
}

// WriteToolchain writes a toolchain inspection in the requested format. In JSON mode it writes
// the full machine envelope identified by metadata (used by doctor and install).
func WriteToolchain(
	writer io.Writer,
	metadata EnvelopeMetadata,
	format OutputFormat,
	result ToolchainResult,
) (allValid bool, err error) {
	switch format {
	case FormatText:
		return writeToolchainText(writer, result)

	case FormatJSON:
		return writeToolchainJSON(writer, metadata, result)

	default:
		return false, fmt.Errorf("unsupported output format %q", format)
	}
}
