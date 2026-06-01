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
	return runExecution(rule.ID, rule.Check, rule.CheckToolIDs(), context, toolStatuses, executors)
}

func RunFix(
	rule contract.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	executors ExecutorRegistry,
) (result contract.ExecutionResult, err error) {
	return runExecution(rule.ID, rule.Fix, rule.FixToolIDs(), context, toolStatuses, executors)
}

func runExecution(
	ruleID string,
	execution contract.ExecutionSpec,
	toolIDs []string,
	context Context,
	toolStatuses map[string]toolchain.Status,
	executors ExecutorRegistry,
) (result contract.ExecutionResult, err error) {
	if execution.Empty() {
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

	executor, found := executors[execution.Kind]
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf(
			"rule %s uses unknown executor %q",
			ruleID,
			execution.Executor(),
		)
	}

	return executor(context, execution, toolStatuses)
}
