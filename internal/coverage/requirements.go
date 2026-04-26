package coverage

import (
	"sort"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/styleguide"
)

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
