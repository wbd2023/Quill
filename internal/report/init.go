package report

import (
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/engine"
)

// WriteInit writes the human-facing result of a completed initialization.
func WriteInit(writer io.Writer, result engine.InitResult) (err error) {
	if _, err = fmt.Fprintf(
		writer,
		"Initialised Quill in %s (preset %s)\n",
		result.Root,
		result.Preset,
	); err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "  wrote STYLE.md, quill.toml")
	return err
}
