package profile

import (
	"fmt"

	"github.com/wbd2023/quill/internal/style"
)

/* -------------------------------- Runtime Binding Completeness -------------------------------- */

// ValidateRuntimeBindings reports an error if any active execution in the effective plan references
// a runtime binding no Pack registered, or if a Fix uses an execution the fix driver set does not
// support. Every active execution specification must resolve to exactly one compatible binding;
// fixes are validated against the fix driver set, which supports file-command and target-command
// executions only.
//
// Tool IDs remain global and file interpreters are keyed by Tool ID; repository scanners, Profile
// checks, target commands, and target checks are keyed by Pack-qualified identity. File-command
// fixes run the tool without interpreting output, so they require no file interpreter.
func ValidateRuntimeBindings(
	effective EffectiveProfile,
	bindings style.RuntimeBindings,
) (err error) {
	for _, rule := range effective.Effective.Rules {
		if err = validateCheckBindings(rule.ID, rule.Check, bindings); err != nil {
			return err
		}

		if err = validateFixBindings(rule.ID, rule.Fix, bindings); err != nil {
			return err
		}
	}

	return nil
}

// validateCheckBindings resolves a rule's Check job against the full check driver set. Every Check
// execution that owns a runtime identity must resolve to exactly one registered binding.
func validateCheckBindings(
	ruleID string,
	job style.Job,
	bindings style.RuntimeBindings,
) (err error) {
	if job == nil {
		return nil
	}

	switch detail := job.(type) {
	case style.ToolchainExecution:
		// Toolchain health is intrinsic to the engine and has no Pack-owned runtime binding.
		return nil

	case style.ExternalCheckJob:
		// External checks are self-describing: the bound job carries its executable and runtime
		// limits, so it needs no Pack-qualified binding registered against the runtime map.
		return nil

	case style.ProfileExecution:
		if !bindings.HasProfileCheck(detail.PackID, detail.Check) {
			return missingBinding(ruleID, "check", "profile check", detail.PackID, detail.Check)
		}

	case style.FileCommandExecution:
		if !bindings.HasFileInterpreter(detail.ToolID) {
			return missingBinding(ruleID, "check", "file interpreter", "", detail.ToolID)
		}

	case style.RepositoryScanExecution:
		if !bindings.HasRepositoryScanner(detail.PackID, detail.Scanner) {
			return missingBinding(ruleID, "check", "repository scanner",
				detail.PackID, detail.Scanner)
		}

	case style.TargetCommandJob:
		if !bindings.HasTargetCommand(detail.PackID, detail.Language, detail.Action) {
			return missingBinding(ruleID, "check", "target command", detail.PackID,
				detail.Language+"/"+detail.Action)
		}

	case style.TargetCheckJob:
		if !bindings.HasTargetCheck(detail.PackID, detail.Language, detail.Check) {
			return missingBinding(ruleID, "check", "target check", detail.PackID,
				detail.Language+"/"+detail.Check)
		}

	default:
		return fmt.Errorf("rule %q check uses unknown execution job %T", ruleID, job)
	}

	return nil
}

// validateFixBindings resolves a rule's Fix job against the fix driver set, which supports
// file-command and target-command executions only. A file-command fix runs the tool without
// interpreting output, so it requires no file interpreter; a target-command fix still resolves its
// Pack-qualified command binding. Any other execution is an unsupported fix and is rejected at
// preparation.
func validateFixBindings(
	ruleID string,
	job style.Job,
	bindings style.RuntimeBindings,
) (err error) {
	if job == nil {
		return nil
	}

	switch detail := job.(type) {
	case style.FileCommandExecution:
		// Fixes run the tool and succeed or fail on exit code; they never interpret findings, so
		// no file interpreter binding is required.
		return nil

	case style.TargetCommandJob:
		if !bindings.HasTargetCommand(detail.PackID, detail.Language, detail.Action) {
			return missingBinding(ruleID, "fix", "target command", detail.PackID,
				detail.Language+"/"+detail.Action)
		}

	case style.ToolchainExecution,
		style.ProfileExecution,
		style.RepositoryScanExecution,
		style.TargetCheckJob,
		style.ExternalCheckJob:
		return fmt.Errorf(
			"rule %q fix uses unsupported execution %T; "+
				"fixes support file-command and target-command only",
			ruleID,
			job,
		)

	default:
		return fmt.Errorf("rule %q fix uses unknown execution job %T", ruleID, job)
	}

	return nil
}

func missingBinding(
	ruleID string,
	side string,
	kind string,
	packID string,
	local string,
) (err error) {
	if packID == "" {
		return fmt.Errorf(
			"rule %q %s references unregistered %s %q",
			ruleID,
			side,
			kind,
			local,
		)
	}

	return fmt.Errorf(
		"rule %q %s references unregistered %s %q in pack %q",
		ruleID,
		side,
		kind,
		local,
		packID,
	)
}
