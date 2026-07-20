package coverage

import "github.com/wbd2023/Quill/internal/styleguide"

func buildSectionCoverage(
	headings []styleguide.Heading,
	requirements []Requirement,
) (sections []Section) {
	requirementsBySection := make(map[string][]Requirement)
	for _, requirement := range requirements {
		requirementsBySection[requirement.Section] = append(
			requirementsBySection[requirement.Section],
			requirement,
		)
	}

	sections = make([]Section, 0, len(headings))
	for _, heading := range headings {
		entry := Section{
			Section: heading.Section,
			Title:   heading.Title,
		}

		for _, requirement := range requirementsBySection[heading.Section] {
			entry.RequirementCount++

			switch requirement.Mode {
			case ModeAutomated:
				entry.AutomatedCount++
			case ModeReviewOnly:
				entry.ReviewOnlyCount++
			case ModeManualDeferred:
				entry.ManualDeferredCount++
			}
		}

		entry.Status = deriveCoverageStatus(entry)
		sections = append(sections, entry)
	}

	return sections
}

func deriveCoverageStatus(entry Section) (status Status) {
	switch {
	case entry.RequirementCount == 0:
		return StatusManual

	case entry.AutomatedCount == entry.RequirementCount:
		return StatusAutomated

	case entry.AutomatedCount > 0:
		return StatusPartial

	case entry.ReviewOnlyCount == entry.RequirementCount:
		return StatusReviewOnly

	default:
		return StatusManual
	}
}
