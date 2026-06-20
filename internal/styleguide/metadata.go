package styleguide

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/style"
)

/* ----------------------------------------- Definitions ---------------------------------------- */

// metadata constants.
const (
	htmlCommentPrefix = "<!--"
	htmlCommentSuffix = "-->"
	metadataPrefix    = "style:"
)

// metadata constants.
const (
	metadataFieldID     metadataField = "id"
	metadataFieldMode   metadataField = "mode"
	metadataFieldReason metadataField = "reason"
)

const metadataModeReviewOnly = "review_only"

type requirementMetadata struct {
	id     style.RequirementID
	review Review
	source position
}

type metadataFields struct {
	id     string
	mode   string
	reason string
}

type metadataField string

/* ------------------------------------------ Comments ------------------------------------------ */

func parseMetadataComment(source string) (fields metadataFields, found bool, err error) {
	payload, found, err := extractMetadataPayload(source)
	if !found || err != nil {
		return metadataFields{}, found, err
	}

	fields, err = parseMetadataFields(payload)
	if err != nil {
		return metadataFields{}, false, err
	}

	return fields, true, nil
}

func extractMetadataPayload(source string) (payload string, found bool, err error) {
	source = strings.TrimSpace(source)
	body, found := extractCommentBody(source)
	if !found {
		body = strings.TrimSpace(strings.TrimPrefix(source, htmlCommentPrefix))
		if strings.HasPrefix(source, metadataPrefix) || strings.HasPrefix(body, metadataPrefix) {
			return "", false, malformedMetadata(source)
		}

		return "", false, nil
	}

	payload, found = extractStylePayload(body)
	return payload, found, nil
}

func extractCommentBody(source string) (body string, found bool) {
	body, found = strings.CutPrefix(strings.TrimSpace(source), htmlCommentPrefix)
	if !found {
		return "", false
	}

	body, found = strings.CutSuffix(body, htmlCommentSuffix)
	if !found {
		return "", false
	}

	return strings.TrimSpace(body), true
}

func extractStylePayload(body string) (payload string, found bool) {
	payload, found = strings.CutPrefix(body, metadataPrefix)
	if !found {
		return "", false
	}

	return strings.TrimSpace(payload), true
}

/* ------------------------------------------- Fields ------------------------------------------- */

func parseMetadataFields(payload string) (fields metadataFields, err error) {
	if strings.Contains(payload, "\n") {
		return parseBlockMetadata(payload)
	}

	return parseInlineMetadata(payload)
}

func parseMetadataField(name string) (field metadataField, found bool) {
	field = metadataField(strings.TrimSpace(name))
	switch field {
	case metadataFieldID, metadataFieldMode, metadataFieldReason:
		return field, true
	default:
		return "", false
	}
}

func (fields *metadataFields) setField(field metadataField, value string) (err error) {
	if strings.TrimSpace(value) == "" {
		return emptyMetadataField(field)
	}

	switch field {
	case metadataFieldID:
		fields.id = value

	case metadataFieldMode:
		fields.mode = value

	case metadataFieldReason:
		fields.reason = value

	default:
		return unknownMetadataField(field)
	}

	return nil
}

/* ------------------------------------ Requirement Metadata ------------------------------------ */

func buildRequirementMetadata(
	fields metadataFields,
	scheme style.IDScheme,
) (metadata requirementMetadata, err error) {
	id, err := style.ParseRequirementID(fields.id, scheme)
	if err != nil {
		return requirementMetadata{}, fmt.Errorf(
			"invalid requirement id in style metadata comment: %w",
			err,
		)
	}

	hasMode, hasReason := fields.mode != "", fields.reason != ""
	if hasMode != hasReason {
		return requirementMetadata{}, fmt.Errorf(
			"style metadata mode and reason must appear together",
		)
	}

	if hasMode && fields.mode != metadataModeReviewOnly {
		return requirementMetadata{}, fmt.Errorf(
			"unsupported style metadata mode %q",
			fields.mode,
		)
	}

	return requirementMetadata{
		id: id,
		review: Review{
			Only:   hasMode,
			Reason: fields.reason,
		},
	}, nil
}

/* ------------------------------------------- Errors ------------------------------------------- */

func malformedMetadata(source string) (err error) {
	return fmt.Errorf("malformed style metadata comment: %q", source)
}

func emptyMetadataField(field metadataField) (err error) {
	return fmt.Errorf("style metadata field %q must not be empty", field)
}

func unknownMetadataField(field metadataField) (err error) {
	return fmt.Errorf("unknown style metadata field %q", field)
}
