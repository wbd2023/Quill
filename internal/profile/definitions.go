package profile

import (
	"fmt"

	"ciphera/tools/internal/style"
)

func indexRuleDefinitions(
	definitions []style.RuleDefinition,
	availableTools map[string]bool,
) (availableRules map[string]style.RuleDefinition, err error) {
	availableRules = make(map[string]style.RuleDefinition, len(definitions))
	for _, definition := range definitions {
		if err = validateRuleDefinition(definition, availableTools); err != nil {
			return nil, err
		}

		if _, found := availableRules[definition.ID]; found {
			return nil, fmt.Errorf("duplicate rule definition %q", definition.ID)
		}

		availableRules[definition.ID] = definition
	}

	return availableRules, nil
}

func validateRuleDefinition(
	definition style.RuleDefinition,
	availableTools map[string]bool,
) (err error) {
	if isBlank(definition.ID) {
		return fmt.Errorf("rule definition has an empty id")
	}

	if isBlank(definition.Name) {
		return fmt.Errorf("rule definition %q has an empty name", definition.ID)
	}

	if isBlank(string(definition.Group)) {
		return fmt.Errorf("rule definition %q has an empty group", definition.ID)
	}

	if definition.Check == nil {
		return fmt.Errorf("rule definition %q must define check execution", definition.ID)
	}

	if err = validateRuleExecution(
		definition.ID,
		"check",
		definition.Check,
		availableTools,
	); err != nil {
		return err
	}

	if definition.Fix != nil {
		return validateRuleExecution(definition.ID, "fix", definition.Fix, availableTools)
	}

	return nil
}
