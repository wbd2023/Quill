package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/wbd2023/quill/internal/report"
)

func (tool Tool) writeUsageError(usage string, err error) {
	if err != nil {
		_, _ = fmt.Fprintln(tool.stderr, err)
		_, _ = fmt.Fprintln(tool.stderr)
	}

	_, _ = io.WriteString(tool.stderr, usage)
}

func (tool Tool) writeError(err error) {
	_, _ = fmt.Fprintln(tool.stderr, err)
}

func (tool Tool) writeCommandOutput(output string) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return
	}

	_, _ = fmt.Fprintln(tool.stderr, trimmed)
}

// reportCommandError renders an executed command's runtime error. In machine mode it writes a
// stable error envelope to stdout. Only a cancelled operation context receives the cancelled code;
// other errors, including child command timeouts, are operation_failed. It always returns exit code
// 1. In text mode it writes the error to stderr.
func (tool Tool) reportCommandError(
	operationContext context.Context,
	command string,
	format report.OutputFormat,
	err error,
) (exitCode int) {
	if format == report.FormatJSON {
		code := report.ErrorCodeOperationFailed
		if operationContext.Err() != nil {
			code = report.ErrorCodeCancelled
		}
		tool.writeMachineErrorEnvelope(command, code, err)
		return 1
	}

	tool.writeError(err)
	return 1
}

// writeMachineErrorEnvelope writes one error envelope to stdout. It must be the only stdout
// content in machine mode.
func (tool Tool) writeMachineErrorEnvelope(command string, code string, err error) {
	if encodeErr := report.WriteEnvelope(
		tool.stdout, report.NewErrorEnvelope(command, code, err.Error()),
	); encodeErr != nil {
		_, _ = fmt.Fprintln(tool.stderr, encodeErr)
	}
}

// renderToolchainStatus writes a toolchain inspection in the requested format. In JSON mode it
// writes the full machine envelope tagged with command (shared by doctor and install).
func renderToolchainStatus(
	writer io.Writer,
	command string,
	format report.OutputFormat,
	result report.ToolchainResult,
) (allValid bool, err error) {
	view := report.NewToolchainView(result)
	return report.WriteToolchain(writer, command, format, view)
}
