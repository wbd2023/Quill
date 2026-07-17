package execution

import (
	"context"
	"errors"
	"fmt"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------- Errors ------------------------------------------- */

var errRuleBlocked = errors.New("rule blocked by toolchain")

// Executor executes one rule's check or fix job against the repository.
type Executor func(
	ctx context.Context,
	run RunContext,
	job style.Job,
	toolStatuses toolchain.StatusMap,
) (result style.ExecutionResult, err error)

// ExecutorSet holds one executor per execution job type. Fields that are nil are treated as
// "no driver" for this job and produce an empty result.
type ExecutorSet struct {
	Toolchain      Executor
	Profile        Executor
	FileCommand    Executor
	TargetCommand  Executor
	TargetCheck    Executor
	RepositoryScan Executor
}

// IsBlocked reports whether the error indicates a rule was blocked by toolchain health.
func IsBlocked(err error) (blocked bool) {
	return errors.Is(err, errRuleBlocked)
}

// RunRule executes a rule's check against the repository.
func RunRule(
	ctx context.Context,
	rule style.Rule,
	run RunContext,
	toolStatuses toolchain.StatusMap,
	drivers ExecutorSet,
) (result style.ExecutionResult, err error) {
	return runExecution(ctx, rule.ID, rule.Check, rule.CheckToolIDs(), run, toolStatuses, drivers)
}

// RunFix executes a rule's fix against the repository.
func RunFix(
	ctx context.Context,
	rule style.Rule,
	run RunContext,
	toolStatuses toolchain.StatusMap,
	drivers ExecutorSet,
) (result style.ExecutionResult, err error) {
	return runExecution(ctx, rule.ID, rule.Fix, rule.FixToolIDs(), run, toolStatuses, drivers)
}

func runExecution(
	ctx context.Context,
	ruleID string,
	job style.Job,
	toolIDs []string,
	run RunContext,
	toolStatuses toolchain.StatusMap,
	drivers ExecutorSet,
) (result style.ExecutionResult, err error) {
	if job == nil {
		return style.ExecutionResult{}, nil
	}

	if len(toolIDs) > 0 && !toolStatuses.AreAllValid(toolIDs) {
		return style.ExecutionResult{
			Diagnostics: []style.Diagnostic{
				{
					Code:    "toolchain/blocked",
					Message: toolStatuses.ExplainIssues(toolIDs),
				},
			},
		}, errRuleBlocked
	}

	driver, err := driverFor(job, drivers)
	if err != nil {
		return style.ExecutionResult{}, fmt.Errorf("rule %s: %w", ruleID, err)
	}

	if driver == nil {
		return style.ExecutionResult{}, fmt.Errorf(
			"rule %s: no driver registered for execution job %T",
			ruleID,
			job,
		)
	}

	return driver(ctx, run, job, toolStatuses)
}

func driverFor(job style.Job, drivers ExecutorSet) (driver Executor, err error) {
	switch job.(type) {

	case style.ToolchainExecution:
		return drivers.Toolchain, nil

	case style.ProfileExecution:
		return drivers.Profile, nil

	case style.FileCommandExecution:
		return drivers.FileCommand, nil

	case style.TargetCommandJob:
		return drivers.TargetCommand, nil

	case style.TargetCheckJob:
		return drivers.TargetCheck, nil

	case style.RepositoryScanExecution:
		return drivers.RepositoryScan, nil

	default:
		return nil, fmt.Errorf("unknown execution job type %T", job)
	}
}
