package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/wbd2023/Quill/internal/report"
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

func renderToolchainStatus(
	writer io.Writer,
	format report.OutputFormat,
	result report.ToolchainResult,
) (allValid bool, err error) {
	view := report.NewToolchainView(result)
	return report.WriteToolchain(writer, format, view)
}
