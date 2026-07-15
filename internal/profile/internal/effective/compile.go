package effective

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

/* ----------------------------------------- Compilation ---------------------------------------- */

// Compile resolves a validated style profile against available rule and tool definitions.
func Compile(
	config policy.Config,
	definitions style.Definitions,
) (effective style.EffectiveConfig, err error) {
	availableTools := indexToolIDs(definitions.ToolIDs)

	if err = validatePins(config, availableTools); err != nil {
		return style.EffectiveConfig{}, err
	}

	availableRules, err := indexRuleDefinitions(definitions.Rules, availableTools)
	if err != nil {
		return style.EffectiveConfig{}, err
	}

	rules, err := resolveRules(config, availableRules)
	if err != nil {
		return style.EffectiveConfig{}, err
	}

	return style.EffectiveConfig{
		Rules: rules,
	}, nil
}

/* --------------------------------------- Rule Resolution -------------------------------------- */

func resolveRules(
	config policy.Config,
	availableRules map[string]style.RuleDefinition,
) (rules []style.Rule, err error) {
	rules = make([]style.Rule, 0, len(config.Rules))
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
	definition style.RuleDefinition,
) (rule style.Rule, err error) {
	check, err := resolveExecution(config, binding, definition.Check)
	if err != nil {
		return style.Rule{}, err
	}

	fix, err := resolveExecution(config, binding, definition.Fix)
	if err != nil {
		return style.Rule{}, err
	}

	return style.Rule{
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
	execution style.ExecutionSpec,
) (resolved style.ExecutionSpec, err error) {
	if err := validateRuleExecutionBinding(config, binding, execution); err != nil {
		return style.ExecutionSpec{}, err
	}

	targets, err := resolveTargets(config, binding, execution)
	if err != nil {
		return style.ExecutionSpec{}, err
	}

	if execution.UsesTargets() {
		return execution.WithTargets(targets), nil
	}

	return execution, nil
}

func isBlank(value string) (blank bool) {
	return strings.TrimSpace(value) == ""
}
