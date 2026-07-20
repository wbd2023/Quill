package engine

import (
	"context"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

/* ---------------------------------------- Rule Checking --------------------------------------- */

// CheckOptions controls one check operation.
type CheckOptions struct {
	// Scope selects the repository scope. An empty scope uses the repository default.
	Scope style.Scope

	// Mode controls whether all rules or only required rules are selected. The zero value is
	// treated as style.CheckModeAll.
	Mode style.CheckMode

	// StrictRecommendations causes failed recommendations to receive failing check status
	// rather than non-failing recommendation treatment.
	StrictRecommendations bool
}

// CheckResult contains the complete non-presentation result of a check.
type CheckResult struct {
	Scope     style.Scope
	Toolchain ToolchainInspection
	Rules     []RuleCheckResult
}

// RuleCheckResult contains the outcome for one selected rule.
type RuleCheckResult struct {
	Rule           style.Rule
	Status         style.CheckStatus
	Execution      style.ExecutionResult
	ExecutionError error
}

// Check loads the repository, resolves and compiles its profile, selects rules, inspects required
// tools, and executes the selected rules.
//
// Rule execution failures are recorded in RuleCheckResult. The returned error represents
// preparation, cancellation, or a fatal orchestration failure. A partial result may accompany a
// non-nil error.
func (engine *Engine) Check(
	operationContext context.Context,
	options CheckOptions,
) (result CheckResult, operationError error) {
	context, environment, err := engine.prepareRunnerContext(operationContext, options.Scope)
	if err != nil {
		return CheckResult{}, err
	}

	result.Scope = context.Scope

	selected := selectRulesForCheck(context.Effective.Rules, context, options.Mode)
	toolIDs := execution.ToolIDsForRules(selected)
	result.Toolchain = engine.inspectTools(operationContext, context.Tools, toolIDs,
		context.ToolEnvironment)
	toolStatuses := toolchain.NewStatusMap(result.Toolchain.Statuses)

	result.Rules = make([]RuleCheckResult, 0, len(selected))
	for _, rule := range selected {
		executionResult, executionError := execution.RunRule(
			operationContext,
			rule,
			context,
			toolStatuses,
			environment.CheckDrivers,
		)
		result.Rules = append(result.Rules, RuleCheckResult{
			Rule: rule,
			Status: execution.CheckStatus(
				rule, executionResult, executionError, options.StrictRecommendations,
			),
			Execution:      executionResult,
			ExecutionError: executionError,
		})
	}

	return result, nil
}

func selectRulesForCheck(
	available []style.Rule,
	context execution.RunContext,
	mode style.CheckMode,
) (rules []style.Rule) {
	for _, rule := range available {
		if !context.Profile.Repository.HasScopeOverlap(context.Scope, rule.Scope) {
			continue
		}

		if mode == style.CheckModeRequired &&
			rule.Enforcement == style.EnforcementRecommendation {
			continue
		}

		rules = append(rules, rule)
	}

	return rules
}
