package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
)

/* ------------------------------------------- Compile ------------------------------------------ */

func (policy Profile) Compile(
	registry rulepack.Registry,
) (effective EffectiveConfig, err error) {
	if err = policy.Validate(); err != nil {
		return EffectiveConfig{}, err
	}

	builtins := registry.Definitions()
	ruleByID := make(map[string]contract.RuleDefinition, len(builtins.Rules))
	for _, builtin := range builtins.Rules {
		ruleByID[builtin.ID] = builtin
	}

	effective = EffectiveConfig{
		Tools: append([]contract.Tool{}, builtins.Tools...),
		Rules: make([]EffectiveRule, 0, len(policy.Rules)),
	}

	for _, binding := range policy.Rules {
		builtin, found := ruleByID[binding.RuleID]
		if !found {
			return EffectiveConfig{}, fmt.Errorf("unknown builtin rule %q", binding.RuleID)
		}

		if err = policy.validateRuleSpec(binding, builtin); err != nil {
			return EffectiveConfig{}, err
		}

		effective.Rules = append(effective.Rules, contract.Rule{
			RuleDefinition: builtin,
			Level:          binding.Level,
			Scope:          binding.Scope,
			RequirementIDs: append([]string{}, binding.RequirementIDs...),
			ConfigRef:      binding.ConfigRef,
		})
	}

	return effective, nil
}

func (policy Profile) validateRuleSpec(
	binding RuleBinding,
	builtin contract.RuleDefinition,
) (err error) {
	if err = policy.validateConfigRef(binding, builtin); err != nil {
		return err
	}

	if err = policy.validateRequiredPathClasses(binding.RuleID, builtin); err != nil {
		return err
	}

	if err = policy.validateExecutionSpec(binding.RuleID, builtin.Spec); err != nil {
		return err
	}

	if err = policy.validateExecutionSpec(binding.RuleID, builtin.FixSpec); err != nil {
		return err
	}

	if builtin.Spec.Executor == contract.ExecutorGoStyle {
		if err = policy.validateGoStyleConfig(binding.RuleID); err != nil {
			return err
		}
	}

	return nil
}

func (policy Profile) validateExecutionSpec(
	ruleID string,
	spec contract.ExecutionSpec,
) (err error) {
	if spec.Executor == "" {
		return nil
	}

	if spec.Backend != "" {
		err = policy.validateLanguageBackend(ruleID, spec.Backend, spec.Language)
		if err != nil {
			return err
		}
	}

	if spec.FileSet != "" {
		if _, found := policy.FileSet(spec.FileSet); !found {
			return fmt.Errorf(
				"rule %q references unknown file set %q",
				ruleID,
				spec.FileSet,
			)
		}
	}

	return nil
}

func (policy Profile) validateConfigRef(
	binding RuleBinding,
	builtin contract.RuleDefinition,
) (err error) {
	if len(builtin.RequiredConfigRefs) == 0 {
		if binding.ConfigRef != "" {
			return fmt.Errorf(
				"rule %q has unexpected config_ref %q",
				binding.RuleID,
				binding.ConfigRef,
			)
		}

		return nil
	}

	for _, configRef := range builtin.RequiredConfigRefs {
		if binding.ConfigRef == configRef {
			return nil
		}
	}

	return fmt.Errorf(
		"rule %q must use config_ref %q",
		binding.RuleID,
		builtin.RequiredConfigRefs[0],
	)
}

func (policy Profile) validateGoStyleConfig(ruleID string) (err error) {
	parameters := policy.Naming.GoParameters
	if len(parameters.SecretNames) == 0 {
		return fmt.Errorf("rule %q requires naming.go_parameters.secret_names", ruleID)
	}

	if len(parameters.ConstructorCategories) == 0 {
		return fmt.Errorf(
			"rule %q requires naming.go_parameters.constructor_categories",
			ruleID,
		)
	}

	seen := make(map[string]bool, len(parameters.ConstructorCategories))
	for _, category := range parameters.ConstructorCategories {
		if category.Name == "" {
			return fmt.Errorf("rule %q has an unnamed constructor category", ruleID)
		}

		if seen[category.Name] {
			return fmt.Errorf(
				"rule %q duplicates constructor category %q",
				ruleID,
				category.Name,
			)
		}

		seen[category.Name] = true
		if len(category.TypeMarkers) > 0 ||
			len(category.ParameterNames) > 0 ||
			category.UsesSecretNames {
			continue
		}

		return fmt.Errorf(
			"rule %q constructor category %q has no matcher",
			ruleID,
			category.Name,
		)
	}

	return nil
}

func (policy Profile) validateRequiredPathClasses(
	ruleID string,
	builtin contract.RuleDefinition,
) (err error) {
	for _, className := range builtin.RequiredPathClasses {
		if len(policy.Paths.Patterns(className)) == 0 {
			return fmt.Errorf("rule %q requires paths.%s", ruleID, className)
		}
	}

	return nil
}

func (policy Profile) validateLanguageBackend(
	ruleID string,
	backendName string,
	language string,
) (err error) {
	if backendName == "" {
		return fmt.Errorf("rule %q must define a language backend", ruleID)
	}

	backend, found := policy.LanguageBackend(backendName)
	if !found {
		return fmt.Errorf("rule %q references unknown language backend %q", ruleID, backendName)
	}

	if language != "" && backend.Language != language {
		return fmt.Errorf(
			"rule %q requires a %s backend, got %q",
			ruleID,
			language,
			backend.Language,
		)
	}

	return nil
}
