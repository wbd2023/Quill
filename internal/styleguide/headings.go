package styleguide

import (
	"strings"

	"ciphera/tools/internal/requirementid"
)

func parseHeading(text string) (heading Heading, found bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Heading{}, false
	}

	section, remainder, found := strings.Cut(text, " ")
	if !found || !requirementid.ValidSection(section) {
		return Heading{}, false
	}

	title := strings.TrimSpace(remainder)
	if title == "" {
		return Heading{}, false
	}

	return Heading{
		Section: section,
		Title:   title,
	}, true
}
