package profile

import (
	"fmt"

	"github.com/wbd2023/quill/internal/style"
)

/* -------------------------------------- Execution Details ------------------------------------- */

func (validator ruleExecutionValidator) validateToolchainCheck(
	execution style.ToolchainCheck,
) (err error) {
	return validator.validateToolReferences(execution.ToolIDs)
}

func (validator ruleExecutionValidator) validateProfileCheck(
	execution style.ProfileCheck,
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

func (validator ruleExecutionValidator) validateFileCommand(
	execution style.FileCommand,
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

func (validator ruleExecutionValidator) validateTargetCommand(
	execution style.TargetCommandTemplate,
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

func (validator ruleExecutionValidator) validateTargetCheck(
	execution style.TargetCheckTemplate,
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

func (validator ruleExecutionValidator) validateRepositoryScan(
	execution style.RepositoryScan,
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

func (validator ruleExecutionValidator) validateExternalCheck(
	execution style.ExternalCheck,
) (err error) {
	if isBlank(execution.CheckID) {
		return fmt.Errorf(
			"rule definition %q %s must define a check id",
			validator.ruleID,
			validator.label,
		)
	}

	return nil
}
