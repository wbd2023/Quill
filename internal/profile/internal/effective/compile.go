package effective

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* ----------------------------------------- Compilation ---------------------------------------- */

// Compile resolves a validated style profile against available rule and tool definitions.
func Compile(
	config policy.Config,
	definitions contract.Definitions,
) (effective contract.EffectiveConfig, err error) {
	availableTools, err := indexToolDefinitions(definitions.Tools)
	if err != nil {
		return contract.EffectiveConfig{}, err
	}

	tools, err := pinTools(config, definitions.Tools, availableTools)
	if err != nil {
		return contract.EffectiveConfig{}, err
	}

	availableRules, err := indexRuleDefinitions(definitions.Rules, availableTools)
	if err != nil {
		return contract.EffectiveConfig{}, err
	}

	rules, err := resolveRules(config, availableRules)
	if err != nil {
		return contract.EffectiveConfig{}, err
	}

	return contract.EffectiveConfig{
		Tools: tools,
		Rules: rules,
	}, nil
}

/* --------------------------------------- Rule Resolution -------------------------------------- */

func resolveRules(
	config policy.Config,
	availableRules map[string]contract.RuleDefinition,
) (rules []contract.Rule, err error) {
	rules = make([]contract.Rule, 0, len(config.Rules))
	for _, binding := range config.Rules {
		definition, found := availableRules[binding.RuleID]
		if !found {
			return nil, fmt.Errorf(
				"unknown rule definition %q",
				binding.RuleID,
			)
		}

		rule, err := resolveRule(config, binding, definition)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func resolveRule(
	config policy.Config,
	binding policy.RuleBinding,
	definition contract.RuleDefinition,
) (rule contract.Rule, err error) {
	check, err := resolveExecution(config, binding, definition.Check)
	if err != nil {
		return contract.Rule{}, err
	}

	fix, err := resolveExecution(config, binding, definition.Fix)
	if err != nil {
		return contract.Rule{}, err
	}

	return contract.Rule{
		ID:             definition.ID,
		Name:           definition.Name,
		Group:          definition.Group,
		Enforcement:    binding.Enforcement,
		Scope:          binding.Scope,
		RequirementIDs: append([]string{}, binding.RequirementIDs...),
		Check:          check,
		Fix:            fix,
	}, nil
}

func resolveExecution(
	config policy.Config,
	binding policy.RuleBinding,
	execution contract.ExecutionSpec,
) (resolved contract.ExecutionSpec, err error) {
	if err := validateRuleExecutionBinding(config, binding, execution); err != nil {
		return contract.ExecutionSpec{}, err
	}

	targets, err := resolveTargets(config, binding, execution)
	if err != nil {
		return contract.ExecutionSpec{}, err
	}

	if execution.UsesTargets() {
		return execution.WithTargets(targets), nil
	}

	return execution, nil
}

func isBlank(value string) (blank bool) {
	return strings.TrimSpace(value) == ""
}
