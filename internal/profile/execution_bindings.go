package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

/* ----------------------------------- Rule Execution Bindings ---------------------------------- */

func validateRuleExecutionBinding(
	config policy.Config,
	binding policy.RuleBinding,
	template style.Template,
) (err error) {
	if template == nil {
		return nil
	}

	if fileSet := style.Describe(template).FileSet; !isBlank(fileSet) {
		if _, found := config.FileSets.Lookup(fileSet); !found {
			return fmt.Errorf(
				"rule %q references unknown file set %q",
				binding.RuleID,
				fileSet,
			)
		}
	}

	return nil
}

func resolveTargets(
	config policy.Config,
	binding policy.RuleBinding,
	template style.Template,
) (targets []string, err error) {
	requirements := style.Describe(template)
	if !requirements.NeedsTargets {
		return nil, nil
	}

	return inferTargets(
		config,
		binding.RuleID,
		binding.Scope,
		requirements.TargetLanguage,
		requirements.NeedsCheckPaths,
	)
}

func inferTargets(
	config policy.Config,
	ruleID string,
	scope style.Scope,
	language string,
	requiresCheckPaths bool,
) (targets []string, err error) {
	for _, target := range config.Targets {
		if !isBlank(language) && target.Language != language {
			continue
		}

		if !config.Repository.HasScopeOverlap(scope, target.Scope) {
			continue
		}

		if requiresCheckPaths && len(target.CheckPaths) == 0 {
			return nil, fmt.Errorf(
				"rule %q target %q must define check_paths",
				ruleID,
				target.Name,
			)
		}

		targets = append(targets, target.Name)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf(
			"rule %q has no %s targets for scope %q",
			ruleID,
			language,
			scope,
		)
	}

	return targets, nil
}
