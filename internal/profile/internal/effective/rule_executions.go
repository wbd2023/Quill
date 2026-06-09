package effective

import (
	"fmt"

	"ciphera/tools/internal/style"
)

/* --------------------------------------- Rule Executions -------------------------------------- */

type ruleExecutionValidator struct {
	ruleID string
	label  string
	tools  map[string]style.Tool
}

func validateRuleExecution(
	ruleID string,
	label string,
	execution style.ExecutionSpec,
	tools map[string]style.Tool,
) (err error) {
	validator := ruleExecutionValidator{
		ruleID: ruleID,
		label:  label,
		tools:  tools,
	}
	return validator.validate(execution)
}

func (validator ruleExecutionValidator) validate(execution style.ExecutionSpec) (err error) {
	if execution.Detail == nil {
		return fmt.Errorf("rule definition %q %s is missing", validator.ruleID, validator.label)
	}

	switch detail := execution.Detail.(type) {
	case style.ToolchainExecution:
		return validator.validateToolchainExecution(execution.Kind, detail)

	case style.ProjectExecution:
		return validator.validateProjectExecution(execution.Kind, detail)

	case style.FileCommandExecution:
		return validator.validateFileCommandExecution(execution.Kind, detail)

	case style.TargetCommandExecution:
		return validator.validateTargetCommandExecution(execution.Kind, detail)

	case style.TargetCheckExecution:
		return validator.validateTargetCheckExecution(execution.Kind, detail)

	case style.RepositoryScanExecution:
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
	actual style.ExecutionKind,
	expected style.ExecutionKind,
) (err error) {
	if actual == expected {
		return nil
	}

	if isBlank(string(actual)) {
		return fmt.Errorf(
			"rule definition %q %s must define an execution kind",
			validator.ruleID,
			validator.label,
		)
	}

	return fmt.Errorf(
		"rule definition %q %s uses execution kind %q, expected %q",
		validator.ruleID,
		validator.label,
		actual,
		expected,
	)
}

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
