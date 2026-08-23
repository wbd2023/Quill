package report

import "io"

/* ------------------------------------------ List JSON ----------------------------------------- */

// listPayloadJSON carries exactly the selected list section under result. Empty sections are
// omitted so each list response carries a single key matching its selector.
type listPayloadJSON struct {
	Packs  []ListPack  `json:"packs,omitempty"`
	Rules  []ListRule  `json:"rules,omitempty"`
	Tools  []ListTool  `json:"tools,omitempty"`
	Scopes []ListScope `json:"scopes,omitempty"`
}

func writeListJSON(writer io.Writer, metadata EnvelopeMetadata, result ListResult) (err error) {
	payload := listPayloadJSON{}
	switch result.Selector {
	case ListPacks:
		payload.Packs = result.Packs
	case ListRules:
		payload.Rules = result.Rules
	case ListTools:
		payload.Tools = result.Tools
	case ListScopes:
		payload.Scopes = result.Scopes
	}

	return writeResultEnvelope(writer, metadata, payload)
}
