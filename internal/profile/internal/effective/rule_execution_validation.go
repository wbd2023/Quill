package effective

import (
	"fmt"

	"ciphera/tools/internal/style"
)

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
