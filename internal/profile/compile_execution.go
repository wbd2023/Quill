package profile

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* --------------------------------------- Execution Specs -------------------------------------- */

func validateExecutionSpec(
	config policy.Config,
	binding policy.RuleBinding,
	spec contract.ExecutionSpec,
) (err error) {
	if spec.Empty() {
		return nil
	}

	if spec.Kind == "" {
		return fmt.Errorf("rule %q must define an executor kind", binding.RuleID)
	}

	if err = validateTypedExecutionShape(binding.RuleID, spec); err != nil {
		return err
	}

	if isBackendSpec(spec) {
		if len(binding.Backends) == 0 {
			return fmt.Errorf("rule %q must define at least one backend", binding.RuleID)
		}

		err = validateLanguageBackends(
			config,
			binding.RuleID,
			binding.Backends,
			spec.BackendLanguage(),
		)
		if err != nil {
			return err
		}
	} else if len(binding.Backends) > 0 {
		return fmt.Errorf("rule %q has unexpected backends", binding.RuleID)
	}

	if fileSet := spec.FileSetName(); fileSet != "" {
		if _, found := config.FileSet(fileSet); !found {
			return fmt.Errorf(
				"rule %q references unknown file set %q",
				binding.RuleID,
				fileSet,
			)
		}
	}

	return nil
}

func validateTypedExecutionShape(ruleID string, spec contract.ExecutionSpec) (err error) {
	if spec.Detail == nil {
		return fmt.Errorf("rule %q execution spec is missing", ruleID)
	}

	switch detail := spec.Detail.(type) {
	case contract.ToolchainExecution:
		if len(detail.ToolIDs) == 0 {
			return fmt.Errorf("rule %q toolchain spec must define tool IDs", ruleID)
		}

	case contract.ControlPlaneExecution:
		if detail.Check == "" {
			return fmt.Errorf("rule %q control-plane spec must define a check", ruleID)
		}

	case contract.FileCommandExecution:
		if detail.ToolID == "" {
			return fmt.Errorf("rule %q file-command spec must define a tool ID", ruleID)
		}
		if detail.FileSet == "" {
			return fmt.Errorf("rule %q file-command spec must define a file set", ruleID)
		}
		if (detail.ConfigArgument == "") != (detail.ConfigFile == "") {
			return fmt.Errorf(
				"rule %q file-command config argument and file must appear together",
				ruleID,
			)
		}

	case contract.BackendCommandExecution:
		if err = validateBackendCommandSpec(ruleID, "backend command", detail); err != nil {
			return err
		}
		if detail.Action == "" {
			return fmt.Errorf("rule %q backend command spec must define action", ruleID)
		}

	case contract.BackendCheckExecution:
		if len(detail.ToolIDs) == 0 {
			return fmt.Errorf("rule %q backend check spec must define tool IDs", ruleID)
		}
		if detail.Language == "" {
			return fmt.Errorf("rule %q backend check spec must define language", ruleID)
		}
		if detail.Check == "" {
			return fmt.Errorf("rule %q backend check spec must define a check", ruleID)
		}

	case contract.RepositoryScanExecution:
		if detail.Scanner == "" {
			return fmt.Errorf("rule %q repository scan spec must define a scanner", ruleID)
		}

	default:
		return fmt.Errorf("rule %q uses unknown execution spec", ruleID)
	}

	return nil
}

func validateBackendCommandSpec(
	ruleID string,
	name string,
	spec contract.BackendCommandExecution,
) (err error) {
	if len(spec.ToolIDs) == 0 {
		return fmt.Errorf("rule %q %s spec must define tool IDs", ruleID, name)
	}

	if spec.Language == "" {
		return fmt.Errorf("rule %q %s spec must define language", ruleID, name)
	}

	return nil
}

func validateLanguageBackends(
	config policy.Config,
	ruleID string,
	backends []string,
	language string,
) (err error) {
	seen := make(map[string]bool, len(backends))
	for _, backend := range backends {
		if backend == "" {
			return fmt.Errorf("rule %q has an empty backend", ruleID)
		}

		if seen[backend] {
			return fmt.Errorf("rule %q duplicates backend %q", ruleID, backend)
		}

		seen[backend] = true
		if err = validateLanguageBackend(config, ruleID, backend, language); err != nil {
			return err
		}
	}

	return nil
}

func validateLanguageBackend(
	config policy.Config,
	ruleID string,
	backendName string,
	language string,
) (err error) {
	if backendName == "" {
		return fmt.Errorf("rule %q must define a language backend", ruleID)
	}

	backend, found := config.LanguageBackend(backendName)
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
