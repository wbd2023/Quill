package styleguide

import (
	"strconv"
	"strings"
	"unicode"
)

/* ------------------------------------------ Prefixes ------------------------------------------ */

// metadata_inline constants.
const (
	inlineMetadataIDPrefix     = string(metadataFieldID) + "="
	inlineMetadataModePrefix   = string(metadataFieldMode) + "="
	inlineMetadataReasonPrefix = string(metadataFieldReason) + "="
)

/* ------------------------------------------- Parsing ------------------------------------------ */

func parseInlineMetadata(payload string) (fields metadataFields, err error) {
	id, payload, err := parseInlineMetadataID(payload)
	if err != nil {
		return metadataFields{}, err
	}

	mode, payload, err := parseInlineMetadataMode(payload)
	if err != nil {
		return metadataFields{}, err
	}

	reason, err := parseInlineMetadataReason(payload)
	if err != nil {
		return metadataFields{}, err
	}

	return metadataFields{
		id:     id,
		mode:   mode,
		reason: reason,
	}, nil
}

func parseInlineMetadataID(payload string) (id string, rest string, err error) {
	payload = strings.TrimSpace(payload)
	token, rest := cutInlineToken(payload)

	id, found := strings.CutPrefix(token, inlineMetadataIDPrefix)
	if !found {
		return "", "", malformedMetadata(payload)
	}

	if strings.TrimSpace(id) == "" {
		return "", "", emptyMetadataField(metadataFieldID)
	}

	return id, strings.TrimSpace(rest), nil
}

func parseInlineMetadataMode(payload string) (mode string, rest string, err error) {
	if payload == "" {
		return "", "", nil
	}

	token, rest := cutInlineToken(payload)
	value, found := strings.CutPrefix(token, inlineMetadataModePrefix)
	if !found {
		return "", "", inlineMetadataError(payload)
	}

	if strings.TrimSpace(value) == "" {
		return "", "", emptyMetadataField(metadataFieldMode)
	}

	return value, strings.TrimSpace(rest), nil
}

func parseInlineMetadataReason(payload string) (reason string, err error) {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return "", nil
	}

	text, found := strings.CutPrefix(payload, inlineMetadataReasonPrefix)
	if !found {
		return "", inlineMetadataError(payload)
	}

	reason, err = strconv.Unquote(text)
	if err != nil {
		return "", malformedMetadata(payload)
	}

	if strings.TrimSpace(reason) == "" {
		return "", emptyMetadataField(metadataFieldReason)
	}

	return reason, nil
}

/* ------------------------------------------- Tokens ------------------------------------------- */

func cutInlineToken(payload string) (token string, rest string) {
	index := strings.IndexFunc(payload, unicode.IsSpace)
	if index < 0 {
		return payload, ""
	}

	return payload[:index], strings.TrimSpace(payload[index:])
}

/* ------------------------------------------- Errors ------------------------------------------- */

func inlineMetadataError(payload string) (err error) {
	token, _ := cutInlineToken(strings.TrimSpace(payload))
	name, _, found := strings.Cut(token, "=")
	if !found {
		return malformedMetadata(payload)
	}

	_, found = parseMetadataField(name)
	if found {
		return malformedMetadata(payload)
	}

	return unknownMetadataField(metadataField(strings.TrimSpace(name)))
}
