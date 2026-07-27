package report

import "io"

/* ----------------------------------------- Lock Result ---------------------------------------- */

// LockResult is the presentation result of a lock operation.
type LockResult struct {
	// Path is the absolute path of the written quill.lock.
	Path string
	// ArchiveCount is the number of archive-tool entries written.
	ArchiveCount int
}

// WriteLock writes the machine-mode lock result envelope for command. Text-mode lock output is
// owned by the CLI, so this renderer only emits the JSON envelope.
func WriteLock(writer io.Writer, command string, result LockResult) (err error) {
	return writeResultEnvelope(writer, command, newLockJSON(result))
}
