package report

import (
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/engine"
)

/* ----------------------------------------- Lock Result ---------------------------------------- */

// LockResult is the presentation result of a lock operation.
type LockResult struct {
	// Path is the absolute path of the written quill.lock.
	Path string
	// ArchiveCount is the number of archive-tool entries written.
	ArchiveCount int
}

// NewLockResult converts a completed engine lock into the explicit report result.
func NewLockResult(result engine.LockResult) (lock LockResult) {
	return LockResult{Path: result.Path, ArchiveCount: result.ArchiveCount}
}

// WriteLock writes a lock result in the requested format.
func WriteLock(
	writer io.Writer,
	metadata EnvelopeMetadata,
	format OutputFormat,
	result LockResult,
) (err error) {
	switch format {
	case FormatText:
		_, err = fmt.Fprintf(writer, "Wrote %s (%d tools)\n", result.Path, result.ArchiveCount)
		return err
	case FormatJSON:
		return writeResultEnvelope(writer, metadata, newLockJSON(result))
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}
