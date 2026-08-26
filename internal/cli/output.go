package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

func (r Runner) writeUsageError(usage string, err error) {
	if err != nil {
		_, _ = fmt.Fprintln(r.stderr, err)
		_, _ = fmt.Fprintln(r.stderr)
	}

	_, _ = io.WriteString(r.stderr, usage)
}

func (r Runner) writeError(err error) {
	_, _ = fmt.Fprintln(r.stderr, err)
}

func (r Runner) envelopeMetadata(command string) (metadata report.EnvelopeMetadata) {
	return report.EnvelopeMetadata{Command: command, QuillVersion: r.version}
}

// reportCommandError renders an executed command's error. Machine output distinguishes semantic
// argument errors, cancellation returned by the operation, and operational failures. It always
// returns the corresponding conventional exit status. Text output writes only the error.
func (r Runner) reportCommandError(
	command string,
	format report.OutputFormat,
	err error,
) (exitCode int) {
	if format == report.FormatJSON {
		r.writeMachineErrorEnvelope(command, errorCode(err), err)
		if isArgumentError(err) {
			return usageExitCode
		}
		return 1
	}

	r.writeError(err)
	if isArgumentError(err) {
		return usageExitCode
	}
	return 1
}

func errorCode(err error) (code string) {
	if isArgumentError(err) {
		return report.ErrorCodeInvalidArgument
	}

	if errors.Is(err, context.Canceled) {
		return report.ErrorCodeCancelled
	}
	return report.ErrorCodeOperationFailed
}

func isArgumentError(err error) (isArgument bool) {
	var argumentError *engine.ArgumentError
	return errors.As(err, &argumentError)
}

// writeMachineErrorEnvelope writes one error envelope to stdout. It must be the only stdout
// content in machine mode.
func (r Runner) writeMachineErrorEnvelope(command string, code string, err error) {
	if encodeErr := report.WriteEnvelope(
		r.stdout,
		report.NewErrorEnvelope(r.envelopeMetadata(command), code, err.Error()),
	); encodeErr != nil {
		_, _ = fmt.Fprintln(r.stderr, encodeErr)
	}
}
