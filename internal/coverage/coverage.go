package coverage

import (
	"sort"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/styleguide"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	StatusAutomated  Status = "automated"
	StatusPartial    Status = "partial"
	StatusReviewOnly Status = "review_only"
	StatusManual     Status = "manual"
)

/* -------------------------------------------- Types ------------------------------------------- */

type Status string

type Requirement struct {
	ID      string
	Section string
	Text    string
	Mode    styleguide.VerificationMode
	Reason  string
	RuleIDs []string
}

type Section struct {
	Section             string
	Title               string
	Status              Status
	RequirementCount    int
	AutomatedCount      int
	ReviewOnlyCount     int
	ManualDeferredCount int
}

type Report struct {
	Requirements []Requirement
	Sections     []Section
}

/* ------------------------------------- Coverage Reporting ------------------------------------- */

func Build(document styleguide.Document, rules []contract.Rule) (report Report) {
	requirements := buildRequirements(document.Requirements, ruleIDsByRequirement(rules))
	return Report{
		Requirements: requirements,
		Sections:     buildSectionCoverage(document.Headings, requirements),
	}
}

func buildRequirements(
	documented []styleguide.Requirement,
	ruleIDsByRequirement map[string][]string,
) (requirements []Requirement) {
	requirements = make([]Requirement, 0, len(documented))
	for _, documentedRequirement := range documented {
		ruleIDs := append([]string{}, ruleIDsByRequirement[documentedRequirement.ID]...)
		sort.Strings(ruleIDs)

		mode := styleguide.VerificationManualDeferred
		reason := "No automated rule is registered yet for this requirement."
		if len(ruleIDs) > 0 {
			mode = styleguide.VerificationAutomated
			reason = ""
		}
		if documentedRequirement.Mode != "" {
			mode = documentedRequirement.Mode
			reason = documentedRequirement.Reason
		}

		requirements = append(requirements, Requirement{
			ID:      documentedRequirement.ID,
			Section: documentedRequirement.Section,
			Text:    documentedRequirement.Text,
			Mode:    mode,
			Reason:  reason,
			RuleIDs: ruleIDs,
		})
	}

	return requirements
}

func ruleIDsByRequirement(rules []contract.Rule) (grouped map[string][]string) {
	grouped = make(map[string][]string)
	for _, rule := range rules {
		for _, requirementID := range rule.RequirementIDs {
			grouped[requirementID] = appendUniqueStrings(grouped[requirementID], rule.ID)
		}
	}

	return grouped
}

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
			case styleguide.VerificationAutomated:
				entry.AutomatedCount++
			case styleguide.VerificationReviewOnly:
				entry.ReviewOnlyCount++
			case styleguide.VerificationManualDeferred:
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

func appendUniqueStrings(values []string, extras ...string) (combined []string) {
	seen := make(map[string]bool)
	combined = make([]string, 0, len(values)+len(extras))

	for _, value := range append(append([]string{}, values...), extras...) {
		if value == "" || seen[value] {
			continue
		}

		seen[value] = true
		combined = append(combined, value)
	}

	return combined
}
