package runner

import (
	"errors"
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------- Errors ------------------------------------------- */

var (
	errRuleBlocked   = errors.New("rule blocked by toolchain")
	errRuleViolation = errors.New("rule violations found")
)

type Executor func(
	context Context,
	spec contract.ExecutionSpec,
	toolStatuses map[string]toolchain.Status,
) (result contract.ExecutionResult, err error)

type ExecutorRegistry map[contract.ExecutorKind]Executor

func IsBlocked(err error) (blocked bool) {
	return errors.Is(err, errRuleBlocked)
}

func RunRule(
	rule contract.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	executors ExecutorRegistry,
) (result contract.ExecutionResult, err error) {
	return runRuleSpec(rule.ID, rule.Spec, rule.ToolIDs(), context, toolStatuses, executors)
}

func RunFix(
	rule contract.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	executors ExecutorRegistry,
) (result contract.ExecutionResult, err error) {
	return runRuleSpec(rule.ID, rule.FixSpec, rule.FixToolIDs(), context, toolStatuses, executors)
}

func runRuleSpec(
	ruleID string,
	spec contract.ExecutionSpec,
	toolIDs []string,
	context Context,
	toolStatuses map[string]toolchain.Status,
	executors ExecutorRegistry,
) (result contract.ExecutionResult, err error) {
	if spec.Empty() {
		return contract.ExecutionResult{}, nil
	}

	if len(toolIDs) > 0 && !toolchain.AllToolsValid(toolIDs, toolStatuses) {
		return contract.ExecutionResult{
			Diagnostics: []contract.Diagnostic{
				{
					Code:    "toolchain/blocked",
					Message: toolchain.ExplainToolIssues(toolIDs, toolStatuses),
				},
			},
		}, errRuleBlocked
	}

	executor, found := executors[spec.Kind]
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf(
			"rule %s uses unknown executor %q",
			ruleID,
			spec.Executor(),
		)
	}

	return executor(context, spec, toolStatuses)
}
