package styleguide

import (
	"fmt"
	"regexp"
	"strings"
)

/* -------------------------------------- Metadata Parsing -------------------------------------- */

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
	inlineMetadataPattern := regexp.MustCompile(
		`^id=([a-z0-9.-]+)` +
			`(?:\s+mode=(review_only)\s+reason="([^"]+)")?$`,
	)

	matches := inlineMetadataPattern.FindStringSubmatch(body)
	if len(matches) == 0 {
		return RequirementMetadata{}, false, fmt.Errorf(
			"malformed style metadata comment: %q",
			body,
		)
	}

	return buildMetadata(matches[1], matches[2], matches[3], requirementIDFormat)
}

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
	currentKey := ""
	currentValue := ""

	flush := func() error {
		if currentKey == "" {
			return nil
		}

		if _, exists := fields[currentKey]; exists {
			return fmt.Errorf("duplicate %q in style metadata comment", currentKey)
		}
		fields[currentKey] = strings.TrimSpace(currentValue)
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

			currentKey = key
			currentValue = strings.TrimSpace(value)
			continue
		}

		if currentKey == "" {
			return nil, fmt.Errorf("malformed style metadata comment near %q", line)
		}

		currentValue += " " + line
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
