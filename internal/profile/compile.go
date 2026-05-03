package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* ----------------------------------------- Compilation ---------------------------------------- */

// Compile resolves a style profile against available rule and tool definitions.
func Compile(
	config policy.Config,
	definitions contract.Definitions,
) (effective contract.EffectiveConfig, err error) {
	if err = Validate(config); err != nil {
		return contract.EffectiveConfig{}, err
	}

	tools, err := bindPinnedTools(config, definitions.Tools)
	if err != nil {
		return contract.EffectiveConfig{}, err
	}

	ruleByID := make(map[string]contract.RuleDefinition, len(definitions.Rules))
	for _, definition := range definitions.Rules {
		ruleByID[definition.ID] = definition
	}

	effective = contract.EffectiveConfig{
		Tools: tools,
		Rules: make([]contract.Rule, 0, len(config.Rules)),
	}

	for _, binding := range config.Rules {
		definition, found := ruleByID[binding.RuleID]
		if !found {
			return contract.EffectiveConfig{}, fmt.Errorf(
				"unknown rule definition %q",
				binding.RuleID,
			)
		}

		spec, fixSpec, err := bindRuleSpecs(config, binding, definition)
		if err != nil {
			return contract.EffectiveConfig{}, err
		}

		definition.Spec = spec
		definition.FixSpec = fixSpec
		effective.Rules = append(effective.Rules, contract.Rule{
			ID:                       definition.ID,
			Name:                     definition.Name,
			Group:                    definition.Group,
			Spec:                     definition.Spec,
			FixSpec:                  definition.FixSpec,
			RequiredConfigReferences: append([]string{}, definition.RequiredConfigReferences...),
			Level:                    binding.Level,
			Scope:                    binding.Scope,
			RequirementIDs:           append([]string{}, binding.RequirementIDs...),
			ConfigReference:          binding.ConfigReference,
			PathClasses:              append([]string{}, binding.PathClasses...),
		})
	}

	return effective, nil
}

/* --------------------------------------- Rule Execution --------------------------------------- */

func bindRuleSpecs(
	config policy.Config,
	binding policy.RuleBinding,
	definition contract.RuleDefinition,
) (spec contract.ExecutionSpec, fixSpec contract.ExecutionSpec, err error) {
	if err = validateConfigReference(binding, definition); err != nil {
		return contract.ExecutionSpec{}, contract.ExecutionSpec{}, err
	}

	if err = validatePathClasses(config, binding); err != nil {
		return contract.ExecutionSpec{}, contract.ExecutionSpec{}, err
	}

	if err = validateExecutionSpec(config, binding, definition.Spec); err != nil {
		return contract.ExecutionSpec{}, contract.ExecutionSpec{}, err
	}

	if err = validateExecutionSpec(config, binding, definition.FixSpec); err != nil {
		return contract.ExecutionSpec{}, contract.ExecutionSpec{}, err
	}

	spec = bindBackends(binding, definition.Spec)
	fixSpec = bindBackends(binding, definition.FixSpec)
	return spec, fixSpec, nil
}

/* --------------------------------------- Backend Binding -------------------------------------- */

func bindBackends(
	binding policy.RuleBinding,
	spec contract.ExecutionSpec,
) (bound contract.ExecutionSpec) {
	if isBackendSpec(spec) {
		return spec.WithBackends(binding.Backends)
	}

	return spec
}

func isBackendSpec(spec contract.ExecutionSpec) (found bool) {
	switch spec.Detail.(type) {
	case contract.BackendCommandExecution, contract.BackendCheckExecution:
		return true
	default:
		return false
	}
}
