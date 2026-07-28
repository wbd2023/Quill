package engine

import (
	"context"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

/* ----------------------------------------- Rule Fixing ---------------------------------------- */

// FixOptions controls one fix operation.
type FixOptions struct {
	// Scope selects the repository scope. An empty scope uses the repository default.
	Scope style.Scope
}

// FixResult contains toolchain inspection and attempted rule fixes.
type FixResult struct {
	Scope     style.Scope
	Toolchain ToolchainInspection
	Rules     []RuleFixResult
}

// RuleFixResult contains the outcome for one attempted fixer.
type RuleFixResult struct {
	Rule           style.Rule
	Execution      style.ExecutionResult
	ExecutionError error
}

// Fix loads the repository, selects fixable rules for the scope, inspects their required tools,
// and executes fixers.
//
// No fixer is run if the required toolchain is invalid. Toolchain invalidity is represented in the
// result rather than as a preparation error.
func (engine *Engine) Fix(
	operationContext context.Context,
	options FixOptions,
) (result FixResult, operationError error) {
	runContext, driverSets, err := engine.prepareRun(operationContext, options.Scope)
	if err != nil {
		return FixResult{}, err
	}

	result.Scope = runContext.Scope

	rules := selectRulesForFix(runContext.Effective.Rules, runContext)
	toolIDs := execution.ToolIDsForFixes(rules)
	result.Toolchain = engine.inspectTools(operationContext, runContext.Tools, toolIDs,
		runContext.ToolEnvironment)
	if err := operationContext.Err(); err != nil {
		return result, err
	}

	if !result.Toolchain.AllValid {
		return result, nil
	}

	toolStatuses := toolchain.NewStatusMap(result.Toolchain.Statuses)
	result.Rules = make([]RuleFixResult, 0, len(rules))
	for _, rule := range rules {
		executionResult, executionError := execution.RunFix(
			operationContext,
			rule,
			runContext,
			toolStatuses,
			driverSets.fix,
		)
		result.Rules = append(result.Rules, RuleFixResult{
			Rule:           rule,
			Execution:      executionResult,
			ExecutionError: executionError,
		})

		if err := operationContext.Err(); err != nil {
			return result, err
		}
		if executionError != nil {
			return result, nil
		}
	}

	return result, nil
}

func selectRulesForFix(
	available []style.Rule,
	runContext execution.RunContext,
) (rules []style.Rule) {
	for _, rule := range available {
		if !runContext.Profile.Repository.HasScopeOverlap(runContext.Scope, rule.Scope) {
			continue
		}

		if rule.Fix == nil {
			continue
		}

		rules = append(rules, rule)
	}

	return rules
}
