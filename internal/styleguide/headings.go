package styleguide

import "strings"

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
