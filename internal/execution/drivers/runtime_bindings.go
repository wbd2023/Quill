package drivers

import (
	"fmt"

	"github.com/wbd2023/quill/internal/style"
)

/* -------------------------------- Runtime Binding Completeness -------------------------------- */

// Validate reports an error if any active execution in the plan references a runtime binding no
// Pack registered, or if a Fix uses an execution the fix driver set does not support. Every active
// Job must resolve to exactly one compatible binding; fixes are validated against the fix driver
// set, which supports file-command and target-command Jobs only.
//
// Tool IDs remain global and file interpreters are keyed by Tool ID; repository scanners, Profile
// checks, target commands, and target checks are keyed by Pack-qualified identity. File-command
// fixes run the tool without interpreting output, so they require no file interpreter. Pack
// provenance for each binding identity is taken from the Rule, not the Job.
func (bindings Bindings) Validate(plan style.Plan) (err error) {
	for _, rule := range plan.Rules {
		if err = validateCheckBindings(rule, rule.Check, bindings); err != nil {
			return err
		}

		if err = validateFixBindings(rule, rule.Fix, bindings); err != nil {
			return err
		}
	}

	return nil
}

// validateCheckBindings resolves a rule's Check Job against the full check driver set. Every Check
// Job that owns a runtime identity must resolve to exactly one registered binding.
func validateCheckBindings(rule style.Rule, job style.Job, bindings Bindings) (err error) {
	if job == nil {
		return nil
	}

	switch detail := job.(type) {
	case style.ToolchainCheck:
		// Toolchain health is intrinsic to the engine and has no Pack-owned runtime binding.
		return nil

	case style.ExternalCheck:
		// External checks are self-describing: the Job carries its executable and runtime
		// limits, so it needs no Pack-qualified binding registered against the runtime map.
		return nil

	case style.ProfileCheck:
		if !bindings.HasProfileCheck(rule.PackID, detail.Check) {
			return missingBinding(rule.ID, "check", "profile check", rule.PackID, detail.Check)
		}

	case style.FileCommand:
		if !bindings.HasFileInterpreter(detail.ToolID) {
			return missingBinding(rule.ID, "check", "file interpreter", "", detail.ToolID)
		}

	case style.RepositoryScan:
		if !bindings.HasRepositoryScanner(rule.PackID, detail.Scanner) {
			return missingBinding(
				rule.ID, "check", "repository scanner", rule.PackID, detail.Scanner,
			)
		}

	case style.TargetCommandJob:
		if !bindings.HasTargetCommand(rule.PackID, detail.Language, detail.Action) {
			return missingBinding(rule.ID, "check", "target command", rule.PackID,
				detail.Language+"/"+detail.Action)
		}

	case style.TargetCheckJob:
		if !bindings.HasTargetCheck(rule.PackID, detail.Language, detail.Check) {
			return missingBinding(rule.ID, "check", "target check", rule.PackID,
				detail.Language+"/"+detail.Check)
		}

	default:
		return fmt.Errorf("rule %q check uses unknown execution job %T", rule.ID, job)
	}

	return nil
}

// validateFixBindings resolves a rule's Fix Job against the fix driver set, which supports
// file-command and target-command Jobs only. A file-command fix runs the tool without interpreting
// output, so it requires no file interpreter; a target-command fix still resolves its
// Pack-qualified command binding. Any other Job is an unsupported fix and is rejected at
// preparation.
func validateFixBindings(rule style.Rule, job style.Job, bindings Bindings) (err error) {
	if job == nil {
		return nil
	}

	switch detail := job.(type) {
	case style.FileCommand:
		// Fixes run the tool and succeed or fail on exit code; they never interpret findings, so
		// no file interpreter binding is required.
		return nil

	case style.TargetCommandJob:
		if !bindings.HasTargetCommand(rule.PackID, detail.Language, detail.Action) {
			return missingBinding(rule.ID, "fix", "target command", rule.PackID,
				detail.Language+"/"+detail.Action)
		}

	case style.ToolchainCheck,
		style.ProfileCheck,
		style.RepositoryScan,
		style.TargetCheckJob,
		style.ExternalCheck:
		return fmt.Errorf(
			"rule %q fix uses unsupported execution %T; "+
				"fixes support file-command and target-command only",
			rule.ID,
			job,
		)

	default:
		return fmt.Errorf("rule %q fix uses unknown execution job %T", rule.ID, job)
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
