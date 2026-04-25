package styleguide

import (
	"sort"

	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
)

func Requirements(repoRoot string) (requirements []Requirement, err error) {
	documented, err := readRequirements(repoRoot)
	if err != nil {
		return nil, err
	}

	ruleIDsByRequirement, err := ruleIDsByRequirement(repoRoot)
	if err != nil {
		return nil, err
	}
	requirements = make([]Requirement, 0, len(documented))

	for _, documentedRequirement := range documented {
		ruleIDs := append([]string{}, ruleIDsByRequirement[documentedRequirement.ID]...)
		sort.Strings(ruleIDs)

		mode := VerificationManualDeferred
		reason := "No automated rule is registered yet for this requirement."
		if len(ruleIDs) > 0 {
			mode = VerificationAutomated
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

	return requirements, nil
}

func ruleIDsByRequirement(repoRoot string) (grouped map[string][]string, err error) {
	grouped = make(map[string][]string)

	policy, err := profile.Load(repoRoot)
	if err != nil {
		return nil, err
	}

	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		return nil, err
	}

	effective, err := policy.Compile(registry)
	if err != nil {
		return nil, err
	}

	for _, rule := range effective.Rules {
		for _, requirementID := range rule.RequirementIDs {
			grouped[requirementID] = appendUniqueStrings(grouped[requirementID], rule.ID)
		}
	}

	return grouped, nil
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
