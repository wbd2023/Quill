package styleguide

import "strings"

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	RequirementIDFormatSectionSlug = "section_slug"
	minimumCodeSpanMarkerLength    = 2
)

/* --------------------------------------- Requirement IDs -------------------------------------- */

func RequirementSection(requirementID string) (section string) {
	return requirementSection(requirementID, RequirementIDFormatSectionSlug)
}

func requirementSection(requirementID string, format string) (section string) {
	if format != RequirementIDFormatSectionSlug {
		return ""
	}

	firstDot := strings.IndexByte(requirementID, '.')
	if firstDot < 0 {
		return ""
	}

	secondDot := strings.IndexByte(requirementID[firstDot+1:], '.')
	if secondDot < 0 {
		return ""
	}

	secondDot += firstDot + 1
	section = requirementID[:secondDot]
	if !isSectionID(section) {
		return ""
	}

	return section
}

/* --------------------------------------- Heading Parsing -------------------------------------- */

func parseHeadingText(value string) (section string, title string, found bool) {
	value = trimMarkdownHeadingPrefix(value)
	if value == "" {
		return "", "", false
	}

	section, remainder, found := strings.Cut(strings.TrimSpace(value), " ")
	if !found || !isSectionID(section) {
		return "", "", false
	}

	title = strings.TrimSpace(remainder)
	if title == "" {
		return "", "", false
	}

	return section, title, true
}

func trimMarkdownHeadingPrefix(value string) (trimmed string) {
	trimmed = strings.TrimSpace(value)
	if !strings.HasPrefix(trimmed, "#") {
		return trimmed
	}

	trimmed = strings.TrimLeft(trimmed, "#")
	return strings.TrimSpace(trimmed)
}

/* ------------------------------------- Requirement Parsing ------------------------------------ */

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

/* ---------------------------------------- ID Validation --------------------------------------- */

func isRequirementID(value string, format string) (valid bool) {
	section := requirementSection(value, format)
	if section == "" || !strings.HasPrefix(value, section+".") {
		return false
	}

	slug := strings.TrimPrefix(value, section+".")
	if slug == "" {
		return false
	}

	previousHyphen := false
	for _, character := range slug {
		switch {
		case character >= 'a' && character <= 'z':
			previousHyphen = false

		case character >= '0' && character <= '9':
			previousHyphen = false

		case character == '-':
			if previousHyphen {
				return false
			}
			previousHyphen = true

		default:
			return false
		}
	}

	return !previousHyphen
}

func isSectionID(value string) (valid bool) {
	major, minor, found := strings.Cut(value, ".")
	if !found || major == "" || minor == "" {
		return false
	}

	return digitsOnly(major) && digitsOnly(minor)
}

func digitsOnly(value string) (digits bool) {
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}

	return true
}
