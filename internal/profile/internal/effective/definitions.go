package effective

import (
	"fmt"

	"ciphera/tools/internal/style"
)

/* -------------------------------------- Rule Definitions -------------------------------------- */

func indexRuleDefinitions(
	definitions []style.RuleDefinition,
	availableTools map[string]style.Tool,
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
	availableTools map[string]style.Tool,
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

	if definition.Check.Empty() {
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

	if !definition.Fix.Empty() {
		return validateRuleExecution(definition.ID, "fix", definition.Fix, availableTools)
	}

	return nil
}

/* -------------------------------------- Tool Definitions -------------------------------------- */

func indexToolDefinitions(
	definitions []style.Tool,
) (availableTools map[string]style.Tool, err error) {
	availableTools = make(map[string]style.Tool, len(definitions))
	for _, definition := range definitions {
		if isBlank(definition.ID) {
			return nil, fmt.Errorf("tool definition has an empty id")
		}

		if isBlank(definition.Name) {
			return nil, fmt.Errorf("tool definition %q has an empty name", definition.ID)
		}

		if _, found := availableTools[definition.ID]; found {
			return nil, fmt.Errorf("duplicate tool definition %q", definition.ID)
		}

		availableTools[definition.ID] = definition
	}

	return availableTools, nil
}
