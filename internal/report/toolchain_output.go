package report

import (
	"fmt"
	"io"
)

/* -------------------------------------- Toolchain Output -------------------------------------- */

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

func writeToolchainText(
	writer io.Writer,
	view ToolchainView,
) (allValid bool, err error) {
	if _, err = fmt.Fprintln(writer, "Style toolchain"); err != nil {
		return false, err
	}

	for _, status := range view.Result.Statuses {
		state := "PASS"
		details := status.Version
		if !status.Valid {
			state = "FAIL"
			details = status.Issue
			if status.Version != "" {
				details = fmt.Sprintf("%s (found %s)", status.Issue, status.Version)
			}
		}

		if err = writeAlignedColumns(
			writer,
			"  ["+state+"]",
			status.Tool.Name,
			details,
		); err != nil {
			return false, err
		}
	}

	return view.AllValid, nil
}

func writeToolchainJSON(writer io.Writer, view ToolchainView) (allValid bool, err error) {
	err = writeJSON(writer, struct {
		Toolchain ToolchainView `json:"toolchain"`
	}{Toolchain: view})
	return view.AllValid, err
}
