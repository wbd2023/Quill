package coverage

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/styleguide"
)

func Build(document styleguide.Document, rules []contract.Rule) (report Report) {
	requirements := buildRequirements(document.Requirements, ruleIDsByRequirement(rules))
	return Report{
		Requirements: requirements,
		Sections:     buildSectionCoverage(document.Headings, requirements),
	}
}
