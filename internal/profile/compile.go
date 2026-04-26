package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* ------------------------------------------- Compile ------------------------------------------ */

func Compile(
	config policy.Config,
	definitions contract.Definitions,
) (effective contract.EffectiveConfig, err error) {
	if err = Validate(config); err != nil {
		return contract.EffectiveConfig{}, err
	}

	tools, err := bindToolPins(config, definitions.Tools)
	if err != nil {
		return contract.EffectiveConfig{}, err
	}

	ruleByID := make(map[string]contract.RuleDefinition, len(definitions.Rules))
	for _, builtin := range definitions.Rules {
		ruleByID[builtin.ID] = builtin
	}

	effective = contract.EffectiveConfig{
		Tools: tools,
		Rules: make([]contract.Rule, 0, len(config.Rules)),
	}

	for _, binding := range config.Rules {
		builtin, found := ruleByID[binding.RuleID]
		if !found {
			return contract.EffectiveConfig{}, fmt.Errorf("unknown builtin rule %q", binding.RuleID)
		}

		spec, fixSpec, err := bindRuleSpecs(config, binding, builtin)
		if err != nil {
			return contract.EffectiveConfig{}, err
		}

		builtin.Spec = spec
		builtin.FixSpec = fixSpec
		effective.Rules = append(effective.Rules, contract.Rule{
			ID:                 builtin.ID,
			Name:               builtin.Name,
			Group:              builtin.Group,
			Spec:               builtin.Spec,
			FixSpec:            builtin.FixSpec,
			RequiredConfigRefs: append([]string{}, builtin.RequiredConfigRefs...),
			Level:              binding.Level,
			Scope:              binding.Scope,
			RequirementIDs:     append([]string{}, binding.RequirementIDs...),
			ConfigRef:          binding.ConfigRef,
			PathClasses:        append([]string{}, binding.PathClasses...),
		})
	}

	return effective, nil
}

func bindRuleSpecs(
	config policy.Config,
	binding policy.RuleBinding,
	builtin contract.RuleDefinition,
) (spec contract.ExecutionSpec, fixSpec contract.ExecutionSpec, err error) {
	if err = validateConfigRef(binding, builtin); err != nil {
		return contract.ExecutionSpec{}, contract.ExecutionSpec{}, err
	}

	if err = validatePathClasses(config, binding); err != nil {
		return contract.ExecutionSpec{}, contract.ExecutionSpec{}, err
	}

	if err = validateExecutionSpec(config, binding, builtin.Spec); err != nil {
		return contract.ExecutionSpec{}, contract.ExecutionSpec{}, err
	}

	if err = validateExecutionSpec(config, binding, builtin.FixSpec); err != nil {
		return contract.ExecutionSpec{}, contract.ExecutionSpec{}, err
	}

	spec = bindBackends(binding, builtin.Spec)
	fixSpec = bindBackends(binding, builtin.FixSpec)
	return spec, fixSpec, nil
}

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
