package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/style"
)

type ruleExecutionValidator struct {
	ruleID  string
	label   string
	toolIDs map[string]bool
}

func validateRuleExecution(
	ruleID string,
	label string,
	template style.Template,
	toolIDs map[string]bool,
) (err error) {
	validator := ruleExecutionValidator{
		ruleID:  ruleID,
		label:   label,
		toolIDs: toolIDs,
	}
	return validator.validate(template)
}

func (validator ruleExecutionValidator) validate(template style.Template) (err error) {
	if template == nil {
		return fmt.Errorf("rule definition %q %s is missing", validator.ruleID, validator.label)
	}

	switch detail := template.(type) {
	case style.ToolchainExecution:
		return validator.validateToolchainExecution(detail)

	case style.ProfileExecution:
		return validator.validateProfileExecution(detail)

	case style.FileCommandExecution:
		return validator.validateFileCommandExecution(detail)

	case style.TargetCommandTemplate:
		return validator.validateTargetCommandExecution(detail)

	case style.TargetCheckTemplate:
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
