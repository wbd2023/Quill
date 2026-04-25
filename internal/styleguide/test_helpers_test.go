package styleguide

import "ciphera/tools/internal/profile"

func extractHeadings(contents string) (headings []documentHeading) {
	document, err := compileDocument([]byte(contents), testStyleGuideConfig())
	if err != nil {
		panic(err)
	}

	return document.Headings
}

func extractRequirements(contents string) (requirements []documentRequirement) {
	document, err := compileDocument([]byte(contents), testStyleGuideConfig())
	if err != nil {
		panic(err)
	}

	return document.Requirements
}

func testStyleGuideConfig() (config profile.StyleGuideConfig) {
	return profile.StyleGuideConfig{
		Path:                "STYLE.md",
		RequirementIDFormat: RequirementIDFormatSectionSlug,
	}
}
