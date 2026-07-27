package report

import "io"

/* -------------------------------------- Machine Protocol -------------------------------------- */

// SchemaVersion is the immutable machine JSON protocol version (docs/cli-protocol.md).
const SchemaVersion = 1

// Envelope status values.
const (
	StatusOK    = "ok"
	StatusError = "error"
)

// Stable machine error codes. Callers must branch on code, never on message text.
const (
	// ErrorCodeInvalidArgument covers malformed flags, unexpected positional arguments, and
	// repository-root resolution failures.
	ErrorCodeInvalidArgument = "invalid_argument"
	// ErrorCodeOperationFailed covers profile loading, preparation, execution, and rendering
	// failures that are not cancellation.
	ErrorCodeOperationFailed = "operation_failed"
	// ErrorCodeCancelled covers operations whose supplied context was cancelled.
	ErrorCodeCancelled = "cancelled"
)

// Envelope is the single JSON document every machine-mode command writes to stdout. A success
// response carries Result; an error response carries Error. Exactly one is set.
type Envelope struct {
	SchemaVersion int            `json:"schema_version"`
	Command       string         `json:"command"`
	Status        string         `json:"status"`
	Result        any            `json:"result,omitempty"`
	Error         *EnvelopeError `json:"error,omitempty"`
}

// EnvelopeError is the stable machine error inside an error envelope.
type EnvelopeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewResultEnvelope builds a success envelope for command carrying result.
func NewResultEnvelope(command string, result any) (envelope Envelope) {
	return Envelope{
		SchemaVersion: SchemaVersion,
		Command:       command,
		Status:        StatusOK,
		Result:        result,
	}
}

// NewErrorEnvelope builds an error envelope for command with the given stable code and message.
func NewErrorEnvelope(command string, code string, message string) (envelope Envelope) {
	return Envelope{
		SchemaVersion: SchemaVersion,
		Command:       command,
		Status:        StatusError,
		Error:         &EnvelopeError{Code: code, Message: message},
	}
}

// WriteEnvelope serialises exactly one machine envelope to writer.
func WriteEnvelope(writer io.Writer, envelope Envelope) (err error) {
	return writeJSON(writer, envelope)
}

// writeResultEnvelope serialises a success envelope for command carrying result. Command
// renderers share this so every machine response uses the same envelope shape.
func writeResultEnvelope(writer io.Writer, command string, result any) (err error) {
	return WriteEnvelope(writer, NewResultEnvelope(command, result))
}
