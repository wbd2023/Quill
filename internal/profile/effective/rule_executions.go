package effective

import (
	"fmt"

	"ciphera/tools/internal/contract"
)

/* --------------------------------------- Rule Executions -------------------------------------- */

type ruleExecutionValidator struct {
	ruleID string
	label  string
	tools  map[string]contract.Tool
}

func validateRuleExecution(
	ruleID string,
	label string,
	execution contract.ExecutionSpec,
	tools map[string]contract.Tool,
) (err error) {
	validator := ruleExecutionValidator{
		ruleID: ruleID,
		label:  label,
		tools:  tools,
	}
	return validator.validate(execution)
}

func (validator ruleExecutionValidator) validate(execution contract.ExecutionSpec) (err error) {
	if execution.Detail == nil {
		return fmt.Errorf("rule definition %q %s is missing", validator.ruleID, validator.label)
	}

	switch detail := execution.Detail.(type) {
	case contract.ToolchainExecution:
		return validator.validateToolchainExecution(execution.Kind, detail)

	case contract.ProjectExecution:
		return validator.validateProjectExecution(execution.Kind, detail)

	case contract.FileCommandExecution:
		return validator.validateFileCommandExecution(execution.Kind, detail)

	case contract.TargetCommandExecution:
		return validator.validateTargetCommandExecution(execution.Kind, detail)

	case contract.TargetCheckExecution:
		return validator.validateTargetCheckExecution(execution.Kind, detail)

	case contract.RepositoryScanExecution:
		return validator.validateRepositoryScanExecution(execution.Kind, detail)

	default:
		return fmt.Errorf(
			"rule definition %q %s uses unknown execution detail",
			validator.ruleID,
			validator.label,
		)
	}
}

func (validator ruleExecutionValidator) validateExecutionKind(
	actual contract.ExecutorKind,
	expected contract.ExecutorKind,
) (err error) {
	if actual == expected {
		return nil
	}

	if isBlank(string(actual)) {
		return fmt.Errorf(
			"rule definition %q %s must define an executor kind",
			validator.ruleID,
			validator.label,
		)
	}

	return fmt.Errorf(
		"rule definition %q %s uses executor kind %q, expected %q",
		validator.ruleID,
		validator.label,
		actual,
		expected,
	)
}

/* -------------------------------------- Execution Details ------------------------------------- */

func (validator ruleExecutionValidator) validateToolchainExecution(
	kind contract.ExecutorKind,
	execution contract.ToolchainExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, contract.ExecutorToolchain); err != nil {
		return err
	}

	return validator.validateToolReferences(execution.ToolIDs)
}

func (validator ruleExecutionValidator) validateProjectExecution(
	kind contract.ExecutorKind,
	execution contract.ProjectExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, contract.ExecutorProject); err != nil {
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
	kind contract.ExecutorKind,
	execution contract.FileCommandExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, contract.ExecutorFileCommand); err != nil {
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
	kind contract.ExecutorKind,
	execution contract.TargetCommandExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, contract.ExecutorTargetCommand); err != nil {
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
	kind contract.ExecutorKind,
	execution contract.TargetCheckExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, contract.ExecutorTargetCheck); err != nil {
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
	kind contract.ExecutorKind,
	execution contract.RepositoryScanExecution,
) (err error) {
	if err = validator.validateExecutionKind(kind, contract.ExecutorRepositoryScan); err != nil {
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

/* --------------------------------------- Tool References -------------------------------------- */

func (validator ruleExecutionValidator) validateToolReferences(toolIDs []string) (err error) {
	if len(toolIDs) == 0 {
		return fmt.Errorf(
			"rule definition %q %s must define tool IDs",
			validator.ruleID,
			validator.label,
		)
	}

	seen := make(map[string]bool, len(toolIDs))
	for _, toolID := range toolIDs {
		if err = validator.validateToolReference(toolID); err != nil {
			return err
		}

		if seen[toolID] {
			return fmt.Errorf(
				"rule definition %q %s duplicates tool %q",
				validator.ruleID,
				validator.label,
				toolID,
			)
		}

		seen[toolID] = true
	}

	return nil
}

func (validator ruleExecutionValidator) validateToolReference(toolID string) (err error) {
	if isBlank(toolID) {
		return fmt.Errorf(
			"rule definition %q %s has an empty tool ID",
			validator.ruleID,
			validator.label,
		)
	}

	if _, found := validator.tools[toolID]; !found {
		return fmt.Errorf(
			"rule definition %q %s references unknown tool %q",
			validator.ruleID,
			validator.label,
			toolID,
		)
	}

	return nil
}
