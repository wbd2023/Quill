package styleguide

func extractHeadings(contents string) (headings []Heading) {
	document, err := compileDocument([]byte(contents), testStyleGuideConfig())
	if err != nil {
		panic(err)
	}

	return document.Headings
}

func extractRequirements(contents string) (requirements []Requirement) {
	document, err := compileDocument([]byte(contents), testStyleGuideConfig())
	if err != nil {
		panic(err)
	}

	return document.Requirements
}

func testStyleGuideConfig() (config Config) {
	return Config{
		Path:                "STYLE.md",
		RequirementIDFormat: RequirementIDFormatSectionSlug,
	}
}
