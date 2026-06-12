package styleguide

import (
	"strings"
	"testing"

	"ciphera/tools/internal/requirementid"
)

func testStyleGuideConfig() (config Config) {
	return Config{
		Filename: "STYLE.md",
		IDScheme: requirementid.SectionSlug,
	}
}

func styleDocument(lines ...string) (document string) {
	return strings.Join(lines, "\n") + "\n"
}

func parseDocument(t *testing.T, contents string) (document Document) {
	t.Helper()

	document, err := Parse([]byte(contents), testStyleGuideConfig())
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	return document
}

func parseHeadings(t *testing.T, contents string) (headings []Heading) {
	t.Helper()

	return parseDocument(t, contents).Headings
}

func parseRequirements(t *testing.T, contents string) (requirements []Requirement) {
	t.Helper()

	return parseDocument(t, contents).Requirements
}
