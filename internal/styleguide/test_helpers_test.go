package styleguide

import "testing"

func extractHeadings(t *testing.T, contents string) (headings []Heading) {
	t.Helper()

	document, err := compileDocument([]byte(contents), testStyleGuideConfig())
	if err != nil {
		t.Fatalf("compileDocument: %v", err)
	}

	return document.Headings
}

func extractRequirements(t *testing.T, contents string) (requirements []Requirement) {
	t.Helper()

	document, err := compileDocument([]byte(contents), testStyleGuideConfig())
	if err != nil {
		t.Fatalf("compileDocument: %v", err)
	}

	return document.Requirements
}

func testStyleGuideConfig() (config Config) {
	return Config{
		Path:                "STYLE.md",
		RequirementIDFormat: RequirementIDFormatSectionSlug,
	}
}
