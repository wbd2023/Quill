package coverage

import (
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/styleguide"
)

// Build assembles a coverage report mapping STYLE.md requirements to automated rules.
func Build(document styleguide.Document, rules []style.Rule) (report Report) {
	requirements := buildRequirements(document.Requirements, ruleIDsByRequirement(rules))
	return Report{
		Requirements: requirements,
		Sections:     buildSectionCoverage(document.Headings, requirements),
	}
}
