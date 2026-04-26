package styleguide

import "strings"

const minimumCodeSpanMarkerLength = 2

func parseRequirementText(
	value string,
	pendingRequirementID string,
	format string,
) (requirementID string, text string, found bool) {
	itemBody, found := parseListItemBody(value)
	if !found {
		itemBody = strings.TrimSpace(value)
	}

	if requirementID, remainder, found := extractRequirementMarker(itemBody); found {
		if !isRequirementID(requirementID, format) {
			return "", "", false
		}

		text = strings.TrimSpace(remainder)
		if text == "" {
			return "", "", false
		}

		return requirementID, text, true
	}

	if pendingRequirementID == "" || !isRequirementID(pendingRequirementID, format) {
		return "", "", false
	}

	text = strings.TrimSpace(itemBody)
	if text == "" {
		return "", "", false
	}

	return pendingRequirementID, text, true
}

func extractRequirementMarker(
	itemBody string,
) (requirementID string, remainder string, found bool) {
	switch {
	case strings.HasPrefix(itemBody, "`["):
		endIndex := strings.Index(itemBody, "]`")
		if endIndex <= minimumCodeSpanMarkerLength {
			return "", "", false
		}

		return itemBody[minimumCodeSpanMarkerLength:endIndex], itemBody[endIndex+2:], true

	case strings.HasPrefix(itemBody, "["):
		endIndex := strings.IndexByte(itemBody, ']')
		if endIndex <= 1 {
			return "", "", false
		}

		return itemBody[1:endIndex], itemBody[endIndex+1:], true

	default:
		return "", "", false
	}
}
