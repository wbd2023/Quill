package profile

import "fmt"

func (validator ruleExecutionValidator) validateToolReferences(toolIDs []string) (err error) {
	if len(toolIDs) == 0 {
		return fmt.Errorf(
			"rule definition %q %s must define tool IDs",
			validator.ruleID,
			validator.label,
		)
	}

	seen := make(map[string]bool, len(toolIDs))
	for _, toolID := range toolIDs {
		if err = validator.validateToolReference(toolID); err != nil {
			return err
		}

		if seen[toolID] {
			return fmt.Errorf(
				"rule definition %q %s duplicates tool %q",
				validator.ruleID,
				validator.label,
				toolID,
			)
		}

		seen[toolID] = true
	}

	return nil
}

func (validator ruleExecutionValidator) validateToolReference(toolID string) (err error) {
	if isBlank(toolID) {
		return fmt.Errorf(
			"rule definition %q %s has an empty tool ID",
			validator.ruleID,
			validator.label,
		)
	}

	if _, found := validator.toolIDs[toolID]; !found {
		return fmt.Errorf(
			"rule definition %q %s references unknown tool %q",
			validator.ruleID,
			validator.label,
			toolID,
		)
	}

	return nil
}
