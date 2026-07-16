package engine

import (
	"context"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
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
	context, environment, err := engine.prepareRunnerContext(operationContext, options.Scope)
	if err != nil {
		return FixResult{}, err
	}

	result.Scope = context.Scope

	rules := selectRulesForFix(context.Effective.Rules, context)
	toolIDs := runner.ToolIDsForFixes(rules)
	result.Toolchain = engine.inspectTools(context.Tools, toolIDs, context.ToolEnvironment)

	if !result.Toolchain.AllValid {
		return result, nil
	}

	toolStatuses := toolchain.NewStatusMap(result.Toolchain.Statuses)
	result.Rules = make([]RuleFixResult, 0, len(rules))
	for _, rule := range rules {
		execution, executionError := runner.RunFix(
			rule,
			context,
			toolStatuses,
			environment.FixDrivers,
		)
		result.Rules = append(result.Rules, RuleFixResult{
			Rule:           rule,
			Execution:      execution,
			ExecutionError: executionError,
		})

		if executionError != nil {
			return result, nil
		}
	}

	return result, nil
}

func selectRulesForFix(available []style.Rule, context runner.Context) (rules []style.Rule) {
	for _, rule := range available {
		if !context.Profile.Repository.HasScopeOverlap(context.Scope, rule.Scope) {
			continue
		}

		if rule.Fix == nil {
			continue
		}

		rules = append(rules, rule)
	}

	return rules
}

/* --------------------------------------- Full Inspection -------------------------------------- */

// Inspect loads the repository and inspects every configured tool.
func (engine *Engine) Inspect(
	operationContext context.Context,
) (inspection ToolchainInspection, operationError error) {
	context, _, err := engine.prepareRunnerContext(operationContext, "")
	if err != nil {
		return ToolchainInspection{}, err
	}

	toolIDs := make([]string, 0, len(context.Tools))
	for toolID := range context.Tools {
		toolIDs = append(toolIDs, toolID)
	}

	return engine.inspectTools(context.Tools, toolIDs, context.ToolEnvironment), nil
}
