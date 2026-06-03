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

type Driver func(
	context Context,
	spec contract.ExecutionSpec,
	toolStatuses map[string]toolchain.Status,
) (result contract.ExecutionResult, err error)

type DriverRegistry map[contract.ExecutionKind]Driver

func IsBlocked(err error) (blocked bool) {
	return errors.Is(err, errRuleBlocked)
}

func RunRule(
	rule contract.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverRegistry,
) (result contract.ExecutionResult, err error) {
	return runExecution(rule.ID, rule.Check, rule.CheckToolIDs(), context, toolStatuses, drivers)
}

func RunFix(
	rule contract.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverRegistry,
) (result contract.ExecutionResult, err error) {
	return runExecution(rule.ID, rule.Fix, rule.FixToolIDs(), context, toolStatuses, drivers)
}

func runExecution(
	ruleID string,
	execution contract.ExecutionSpec,
	toolIDs []string,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverRegistry,
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

	driver, found := drivers[execution.Kind]
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf(
			"rule %s uses unknown execution kind %q",
			ruleID,
			string(execution.Kind),
		)
	}

	return driver(context, execution, toolStatuses)
}
