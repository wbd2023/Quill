package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

/* ----------------------------------------- Compilation ---------------------------------------- */

// compilePlan resolves a validated style profile against available rule and tool definitions.
func compilePlan(
	config policy.Config,
	definitions style.Definitions,
) (effective style.Plan, err error) {
	availableTools := indexToolIDs(definitions.ToolIDs)

	if err = validatePins(config, availableTools); err != nil {
		return style.Plan{}, err
	}

	availableRules, err := indexRuleDefinitions(definitions.Rules, availableTools)
	if err != nil {
		return style.Plan{}, err
	}

	rules, err := resolveRules(config, availableRules)
	if err != nil {
		return style.Plan{}, err
	}

	return style.Plan{
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
	template style.Template,
) (resolved style.Job, err error) {
	if template == nil {
		return nil, nil
	}

	if err = validateRuleExecutionBinding(config, binding, template); err != nil {
		return nil, err
	}

	targets, err := resolveTargets(config, binding, template)
	if err != nil {
		return nil, err
	}

	return style.Bind(template, targets), nil
}
