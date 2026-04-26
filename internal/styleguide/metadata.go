package styleguide

import (
	"fmt"
	"strings"
)

/* --------------------------------------- Comment Parsing -------------------------------------- */

func parseMetadataComment(
	line string,
	requirementIDFormat string,
) (metadata RequirementMetadata, found bool, err error) {
	trimmed := strings.TrimSpace(line)
	if !strings.Contains(trimmed, "style:") {
		return RequirementMetadata{}, false, nil
	}

	if !strings.HasPrefix(trimmed, "<!--") || !strings.HasSuffix(trimmed, "-->") {
		return RequirementMetadata{}, false, fmt.Errorf(
			"malformed style metadata comment: %q",
			trimmed,
		)
	}

	body := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "<!--"), "-->"))
	if !strings.HasPrefix(body, "style:") {
		return RequirementMetadata{}, false, fmt.Errorf(
			"malformed style metadata comment: %q",
			body,
		)
	}

	body = strings.TrimSpace(strings.TrimPrefix(body, "style:"))
	if strings.Contains(body, "\n") {
		return parseBlockMetadata(body, requirementIDFormat)
	}

	return parseInlineMetadata(body, requirementIDFormat)
}

func parseInlineMetadata(
	body string,
	requirementIDFormat string,
) (metadata RequirementMetadata, found bool, err error) {
	requirementID, rest, err := parseInlineMetadataID(body)
	if err != nil {
		return RequirementMetadata{}, false, err
	}

	mode, reason, err := parseInlineMetadataReview(rest)
	if err != nil {
		return RequirementMetadata{}, false, err
	}

	return buildMetadata(requirementID, mode, reason, requirementIDFormat)
}

func parseInlineMetadataID(body string) (requirementID string, rest string, err error) {
	body = strings.TrimSpace(body)
	if !strings.HasPrefix(body, "id=") {
		return "", "", malformedMetadata(body)
	}

	requirementID, rest, _ = strings.Cut(strings.TrimPrefix(body, "id="), " ")
	if requirementID == "" {
		return "", "", malformedMetadata(body)
	}

	return requirementID, strings.TrimSpace(rest), nil
}

func parseInlineMetadataReview(rest string) (mode string, reason string, err error) {
	if rest == "" {
		return "", "", nil
	}

	modeField, reasonField, hasReason := strings.Cut(rest, " ")
	mode, found := strings.CutPrefix(modeField, "mode=")
	if !found {
		return "", "", malformedMetadata(rest)
	}

	if !hasReason {
		return mode, "", nil
	}

	reasonSource, found := strings.CutPrefix(strings.TrimSpace(reasonField), `reason="`)
	if !found {
		return "", "", malformedMetadata(rest)
	}

	reason, tail, found := strings.Cut(reasonSource, `"`)
	if !found || tail != "" {
		return "", "", malformedMetadata(rest)
	}

	return mode, reason, nil
}

/* ---------------------------------------- Block Fields ---------------------------------------- */

func parseBlockMetadata(
	body string,
	requirementIDFormat string,
) (metadata RequirementMetadata, found bool, err error) {
	fields, err := parseMetadataFields(body)
	if err != nil {
		return RequirementMetadata{}, false, err
	}

	return buildMetadata(fields["id"], fields["mode"], fields["reason"], requirementIDFormat)
}

func parseMetadataFields(body string) (fields map[string]string, err error) {
	fields = make(map[string]string)
	fieldKey := ""
	fieldText := ""

	flush := func() error {
		if fieldKey == "" {
			return nil
		}

		if _, exists := fields[fieldKey]; exists {
			return fmt.Errorf("duplicate %q in style metadata comment", fieldKey)
		}
		fields[fieldKey] = strings.TrimSpace(fieldText)
		return nil
	}

	for _, rawLine := range strings.Split(body, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		key, value, hasAssignment := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		if hasAssignment && isMetadataFieldKey(key) {
			if err := flush(); err != nil {
				return nil, err
			}

			fieldKey = key
			fieldText = strings.TrimSpace(value)
			continue
		}

		if fieldKey == "" {
			return nil, fmt.Errorf("malformed style metadata comment near %q", line)
		}

		fieldText += " " + line
	}

	if err := flush(); err != nil {
		return nil, err
	}

	return fields, nil
}

func isMetadataFieldKey(value string) (found bool) {
	switch value {
	case "id", "mode", "reason":
		return true
	default:
		return false
	}
}

func malformedMetadata(body string) (err error) {
	return fmt.Errorf("malformed style metadata comment: %q", body)
}

/* --------------------------------------- Metadata Values -------------------------------------- */

func buildMetadata(
	requirementID string,
	mode string,
	reason string,
	requirementIDFormat string,
) (metadata RequirementMetadata, found bool, err error) {
	if !isRequirementID(requirementID, requirementIDFormat) {
		return RequirementMetadata{}, false, fmt.Errorf(
			"invalid requirement id in style metadata comment",
		)
	}

	if (mode == "") != (reason == "") {
		return RequirementMetadata{}, false, fmt.Errorf(
			"style metadata mode and reason must appear together",
		)
	}

	if mode != "" && VerificationMode(mode) != VerificationReviewOnly {
		return RequirementMetadata{}, false, fmt.Errorf(
			"unsupported style metadata mode %q",
			mode,
		)
	}

	return RequirementMetadata{
		ID:     requirementID,
		Mode:   VerificationMode(mode),
		Reason: reason,
	}, true, nil
}
