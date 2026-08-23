package report

import (
	"fmt"
	"io"
)

// WriteInstall writes the post-install toolchain result. In text mode the success line is part of
// the operation presentation; machine mode remains the single shared toolchain envelope.
func WriteInstall(
	writer io.Writer,
	metadata EnvelopeMetadata,
	format OutputFormat,
	result ToolchainResult,
) (allValid bool, err error) {
	allValid, err = WriteToolchain(writer, metadata, format, result)
	if err != nil || format == FormatJSON || !allValid {
		return allValid, err
	}

	_, err = fmt.Fprintln(writer, "Style tools installed.")
	return allValid, err
}
