package execution

import (
	"context"
	"errors"
	"fmt"

	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

/* ------------------------------------------- Errors ------------------------------------------- */

var errRuleBlocked = errors.New("rule blocked by toolchain")

// Driver executes one rule's check or fix Job against the repository. The Job carries the bound
// execution value and resolved targets; the Rule carries Pack and rule provenance so request and
// error attribution derive from it rather than from the job or result.
type Driver func(
	ctx context.Context,
	run RunContext,
	rule style.Rule,
	job style.Job,
	toolStatuses toolchain.StatusMap,
) (result style.ExecutionResult, err error)

// DriverSet holds one Driver per execution Job family. A nil field is an unresolved driver for
// that family: runExecution returns an error rather than silently producing an empty result, so a
// missing binding fails loudly instead of masquerading as a clean check.
type DriverSet struct {
	Toolchain      Driver
	Profile        Driver
	FileCommand    Driver
	TargetCommand  Driver
	TargetCheck    Driver
	RepositoryScan Driver
	ExternalCheck  Driver
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
	drivers DriverSet,
) (result style.ExecutionResult, err error) {
	return runExecution(
		ctx, rule, rule.Check, rule.CheckToolIDs(),
		run, toolStatuses, drivers,
	)
}

// RunFix executes a rule's fix against the repository.
func RunFix(
	ctx context.Context,
	rule style.Rule,
	run RunContext,
	toolStatuses toolchain.StatusMap,
	drivers DriverSet,
) (result style.ExecutionResult, err error) {
	return runExecution(
		ctx, rule, rule.Fix, rule.FixToolIDs(),
		run, toolStatuses, drivers,
	)
}

func runExecution(
	ctx context.Context,
	rule style.Rule,
	job style.Job,
	toolIDs []string,
	run RunContext,
	toolStatuses toolchain.StatusMap,
	drivers DriverSet,
) (result style.ExecutionResult, err error) {
	if job == nil {
		return style.ExecutionResult{}, nil
	}

	if _, isToolchain := job.(style.ToolchainCheck); !isToolchain &&
		len(toolIDs) > 0 && !toolStatuses.AreAllValid(toolIDs) {
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
		return style.ExecutionResult{}, fmt.Errorf("rule %s: %w", rule.ID, err)
	}

	if driver == nil {
		return style.ExecutionResult{}, fmt.Errorf(
			"rule %s: no driver registered for execution job %T",
			rule.ID,
			job,
		)
	}

	return driver(ctx, run, rule, job, toolStatuses)
}

func driverFor(job style.Job, drivers DriverSet) (driver Driver, err error) {
	switch job.(type) {

	case style.ToolchainCheck:
		return drivers.Toolchain, nil

	case style.ProfileCheck:
		return drivers.Profile, nil

	case style.FileCommand:
		return drivers.FileCommand, nil

	case style.TargetCommandJob:
		return drivers.TargetCommand, nil

	case style.TargetCheckJob:
		return drivers.TargetCheck, nil

	case style.RepositoryScan:
		return drivers.RepositoryScan, nil

	case style.ExternalCheck:
		return drivers.ExternalCheck, nil

	default:
		return nil, fmt.Errorf("unknown execution job type %T", job)
	}
}
