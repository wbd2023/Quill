package coverage

import (
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/styleguide"
)

// Build assembles a coverage report mapping STYLE.md requirements to automated rules.
func Build(document styleguide.Document, rules []style.Rule) (report Report) {
	requirements := buildRequirements(document.Requirements, ruleIDsByRequirement(rules))
	return Report{
		Requirements: requirements,
		Sections:     buildSectionCoverage(document.Headings, requirements),
	}
}
