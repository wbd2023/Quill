package report

import "io"

func writeExplainJSON(
	writer io.Writer,
	metadata EnvelopeMetadata,
	result ExplainResult,
) (err error) {
	return writeResultEnvelope(writer, metadata, result)
}
