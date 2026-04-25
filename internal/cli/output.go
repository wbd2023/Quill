package cli

import (
	"fmt"
	"io"
	"strings"

	"ciphera/tools/internal/report"
)

func (tool CLI) writeUsageError(usage string, err error) {
	if err != nil {
		_, _ = fmt.Fprintln(tool.stderr, err)
		_, _ = fmt.Fprintln(tool.stderr)
	}

	_, _ = io.WriteString(tool.stderr, usage)
}

func (tool CLI) writeError(err error) {
	_, _ = fmt.Fprintln(tool.stderr, err)
}

func (tool CLI) writeCommandOutput(output string) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return
	}

	_, _ = fmt.Fprintln(tool.stderr, trimmed)
}

func renderToolchainStatus(
	writer io.Writer,
	format report.OutputFormat,
	result report.ToolchainResult,
) (allValid bool, err error) {
	view := report.NewToolchainView(result)
	return report.WriteToolchain(writer, format, view)
}
