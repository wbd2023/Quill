package effective

import (
	"fmt"

	"ciphera/tools/internal/style"
)

/* -------------------------------------- Execution Details ------------------------------------- */

func (validator ruleExecutionValidator) validateToolchainExecution(
	execution style.ToolchainExecution,
) (err error) {
	return validator.validateToolReferences(execution.ToolIDs)
}

func (validator ruleExecutionValidator) validateProfileExecution(
	execution style.ProfileExecution,
) (err error) {
	if isBlank(execution.Check) {
		return fmt.Errorf(
			"rule definition %q %s must define a check",
			validator.ruleID,
			validator.label,
		)
	}

	return nil
}

func (validator ruleExecutionValidator) validateFileCommandExecution(
	execution style.FileCommandExecution,
) (err error) {
	if isBlank(execution.ToolID) {
		return fmt.Errorf(
			"rule definition %q %s must define a tool ID",
			validator.ruleID,
			validator.label,
		)
	}

	if err = validator.validateToolReference(execution.ToolID); err != nil {
		return err
	}

	if isBlank(execution.FileSet) {
		return fmt.Errorf(
			"rule definition %q %s must define a file set",
			validator.ruleID,
			validator.label,
		)
	}

	if isBlank(execution.ConfigArgument) != isBlank(execution.ConfigFile) {
		return fmt.Errorf(
			"rule definition %q %s config argument and file must appear together",
			validator.ruleID,
			validator.label,
		)
	}

	return nil
}

func (validator ruleExecutionValidator) validateTargetCommandExecution(
	execution style.TargetCommandExecution,
) (err error) {
	if err = validator.validateToolReferences(execution.ToolIDs); err != nil {
		return err
	}

	if isBlank(execution.Language) {
		return fmt.Errorf(
			"rule definition %q %s must define language",
			validator.ruleID,
			validator.label,
		)
	}

	if isBlank(execution.Action) {
		return fmt.Errorf(
			"rule definition %q %s must define an action",
			validator.ruleID,
			validator.label,
		)
	}

	return nil
}

func (validator ruleExecutionValidator) validateTargetCheckExecution(
	execution style.TargetCheckExecution,
) (err error) {
	if err = validator.validateToolReferences(execution.ToolIDs); err != nil {
		return err
	}

	if isBlank(execution.Language) {
		return fmt.Errorf(
			"rule definition %q %s must define language",
			validator.ruleID,
			validator.label,
		)
	}

	if isBlank(execution.Check) {
		return fmt.Errorf(
			"rule definition %q %s must define a check",
			validator.ruleID,
			validator.label,
		)
	}

	return nil
}

func (validator ruleExecutionValidator) validateRepositoryScanExecution(
	execution style.RepositoryScanExecution,
) (err error) {
	if isBlank(execution.Scanner) {
		return fmt.Errorf(
			"rule definition %q %s must define a scanner",
			validator.ruleID,
			validator.label,
		)
	}

	return nil
}
