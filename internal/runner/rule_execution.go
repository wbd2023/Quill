package runner

import (
	"errors"
	"fmt"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------- Errors ------------------------------------------- */

var (
	errRuleBlocked   = errors.New("rule blocked by toolchain")
	errRuleViolation = errors.New("rule violations found")
)

// Driver is driver.
type Driver func(
	context Context,
	spec style.ExecutionSpec,
	toolStatuses map[string]toolchain.Status,
) (result style.ExecutionResult, err error)

// DriverRegistry is driver registry.
type DriverRegistry map[style.ExecutionKind]Driver

// IsBlocked is blocked.
func IsBlocked(err error) (blocked bool) {
	return errors.Is(err, errRuleBlocked)
}

// RunRule run rule.
func RunRule(
	rule style.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverRegistry,
) (result style.ExecutionResult, err error) {
	return runExecution(rule.ID, rule.Check, rule.CheckToolIDs(), context, toolStatuses, drivers)
}

// RunFix run fix.
func RunFix(
	rule style.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverRegistry,
) (result style.ExecutionResult, err error) {
	return runExecution(rule.ID, rule.Fix, rule.FixToolIDs(), context, toolStatuses, drivers)
}

func runExecution(
	ruleID string,
	execution style.ExecutionSpec,
	toolIDs []string,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverRegistry,
) (result style.ExecutionResult, err error) {
	if execution.Empty() {
		return style.ExecutionResult{}, nil
	}

	if len(toolIDs) > 0 && !toolchain.AreAllToolsValid(toolIDs, toolStatuses) {
		return style.ExecutionResult{
			Diagnostics: []style.Diagnostic{
				{
					Code:    "toolchain/blocked",
					Message: toolchain.ExplainToolIssues(toolIDs, toolStatuses),
				},
			},
		}, errRuleBlocked
	}

	driver, found := drivers[execution.Kind]
	if !found {
		return style.ExecutionResult{}, fmt.Errorf(
			"rule %s uses unknown execution kind %q",
			ruleID,
			string(execution.Kind),
		)
	}

	return driver(context, execution, toolStatuses)
}
