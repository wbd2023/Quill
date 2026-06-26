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
		return validator.validateToolchainExecution(detail)

	case style.ProfileExecution:
		return validator.validateProfileExecution(detail)

	case style.FileCommandExecution:
		return validator.validateFileCommandExecution(detail)

	case style.TargetCommandExecution:
		return validator.validateTargetCommandExecution(detail)

	case style.TargetCheckExecution:
		return validator.validateTargetCheckExecution(detail)

	case style.RepositoryScanExecution:
		return validator.validateRepositoryScanExecution(detail)

	default:
		return fmt.Errorf(
			"rule definition %q %s uses unknown execution detail",
			validator.ruleID,
			validator.label,
		)
	}
}
