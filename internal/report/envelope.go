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
	// ErrorCodeInvalidArgument covers malformed grammar and repository-dependent selections that
	// are syntactically valid but unavailable in the loaded profile.
	ErrorCodeInvalidArgument = "invalid_argument"
	// ErrorCodeOperationFailed covers repository discovery, profile loading, preparation,
	// execution, and rendering failures that are not cancellation.
	ErrorCodeOperationFailed = "operation_failed"
	// ErrorCodeCancelled covers operations that return an error wrapping context.Canceled.
	ErrorCodeCancelled = "cancelled"
)

// EnvelopeMetadata identifies the executable and command that produced an envelope.
type EnvelopeMetadata struct {
	Command      string
	QuillVersion string
}

// Envelope is the single JSON document every machine-mode command writes to stdout. A success
// response carries Result; an error response carries Error. Exactly one is set.
type Envelope struct {
	SchemaVersion int            `json:"schema_version"`
	QuillVersion  string         `json:"quill_version"`
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

// NewResultEnvelope builds a success envelope for metadata carrying result.
func NewResultEnvelope(metadata EnvelopeMetadata, result any) (envelope Envelope) {
	return Envelope{
		SchemaVersion: SchemaVersion,
		QuillVersion:  metadata.QuillVersion,
		Command:       metadata.Command,
		Status:        StatusOK,
		Result:        result,
	}
}

// NewErrorEnvelope builds an error envelope for metadata with the given stable code and message.
func NewErrorEnvelope(
	metadata EnvelopeMetadata,
	code string,
	message string,
) (envelope Envelope) {
	return Envelope{
		SchemaVersion: SchemaVersion,
		QuillVersion:  metadata.QuillVersion,
		Command:       metadata.Command,
		Status:        StatusError,
		Error:         &EnvelopeError{Code: code, Message: message},
	}
}

// WriteEnvelope serialises exactly one machine envelope to writer.
func WriteEnvelope(writer io.Writer, envelope Envelope) (err error) {
	return writeJSON(writer, envelope)
}

// writeResultEnvelope serialises a success envelope for metadata carrying result. Command
// renderers share this so every machine response uses the same envelope shape.
func writeResultEnvelope(
	writer io.Writer,
	metadata EnvelopeMetadata,
	result any,
) (err error) {
	return WriteEnvelope(writer, NewResultEnvelope(metadata, result))
}
