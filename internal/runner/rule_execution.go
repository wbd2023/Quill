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

// Driver executes one rule's check or fix against the repository.
type Driver func(
	context Context,
	spec style.ExecutionSpec,
	toolStatuses map[string]toolchain.Status,
) (result style.ExecutionResult, err error)

// DriverSet holds one driver per execution detail type. Fields that are nil are treated as "no
// driver for this execution" and produce an empty result.
type DriverSet struct {
	Toolchain      Driver
	Profile        Driver
	FileCommand    Driver
	TargetCommand  Driver
	TargetCheck    Driver
	RepositoryScan Driver
}

// IsBlocked reports whether the error indicates a rule was blocked by toolchain health.
func IsBlocked(err error) (blocked bool) {
	return errors.Is(err, errRuleBlocked)
}

// RunRule executes a rule's check against the repository.
func RunRule(
	rule style.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverSet,
) (result style.ExecutionResult, err error) {
	return runExecution(rule.ID, rule.Check, rule.CheckToolIDs(), context, toolStatuses, drivers)
}

// RunFix executes a rule's fix against the repository.
func RunFix(
	rule style.Rule,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverSet,
) (result style.ExecutionResult, err error) {
	return runExecution(rule.ID, rule.Fix, rule.FixToolIDs(), context, toolStatuses, drivers)
}

func runExecution(
	ruleID string,
	execution style.ExecutionSpec,
	toolIDs []string,
	context Context,
	toolStatuses map[string]toolchain.Status,
	drivers DriverSet,
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

	driver, err := driverFor(execution.Detail, drivers)
	if err != nil {
		return style.ExecutionResult{}, fmt.Errorf("rule %s: %w", ruleID, err)
	}

	if driver == nil {
		return style.ExecutionResult{}, nil
	}

	return driver(context, execution, toolStatuses)
}

func driverFor(detail style.ExecutionDetail, drivers DriverSet) (driver Driver, err error) {
	switch detail.(type) {

	case style.ToolchainExecution:
		return drivers.Toolchain, nil

	case style.ProfileExecution:
		return drivers.Profile, nil

	case style.FileCommandExecution:
		return drivers.FileCommand, nil

	case style.TargetCommandExecution:
		return drivers.TargetCommand, nil

	case style.TargetCheckExecution:
		return drivers.TargetCheck, nil

	case style.RepositoryScanExecution:
		return drivers.RepositoryScan, nil

	default:
		return nil, fmt.Errorf("unknown execution detail type %T", detail)
	}
}
