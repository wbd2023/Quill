package styleguide

import "strings"

const (
	RequirementIDFormatSectionSlug = "section_slug"
)

/* ------------------------------------ Requirement Sections ------------------------------------ */

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
