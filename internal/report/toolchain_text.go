package report

import (
	"fmt"
	"io"
)

func writeToolchainText(
	writer io.Writer,
	result ToolchainResult,
) (allValid bool, err error) {
	if _, err = fmt.Fprintln(writer, "Style toolchain"); err != nil {
		return false, err
	}

	for _, status := range result.Statuses {
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

	return result.AllValid, nil
}
