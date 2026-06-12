package effective

import (
	"fmt"

	"ciphera/tools/internal/style"
)

/* -------------------------------------- Execution Details ------------------------------------- */

func (validator ruleExecutionValidator) validateToolchainExecution(
	kind style.ExecutionKind,
	execution style.ToolchainExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, style.ExecutionToolchain); err != nil {
		return err
	}

	return validator.validateToolReferences(execution.ToolIDs)
}

func (validator ruleExecutionValidator) validateProjectExecution(
	kind style.ExecutionKind,
	execution style.ProjectExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, style.ExecutionProject); err != nil {
		return err
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

func (validator ruleExecutionValidator) validateFileCommandExecution(
	kind style.ExecutionKind,
	execution style.FileCommandExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, style.ExecutionFileCommand); err != nil {
		return err
	}

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
	kind style.ExecutionKind,
	execution style.TargetCommandExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, style.ExecutionTargetCommand); err != nil {
		return err
	}

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
	kind style.ExecutionKind,
	execution style.TargetCheckExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, style.ExecutionTargetCheck); err != nil {
		return err
	}

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
	kind style.ExecutionKind,
	execution style.RepositoryScanExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, style.ExecutionRepositoryScan); err != nil {
		return err
	}

	if isBlank(execution.Scanner) {
		return fmt.Errorf(
			"rule definition %q %s must define a scanner",
			validator.ruleID,
			validator.label,
		)
	}

	return nil
}
