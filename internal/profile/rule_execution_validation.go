package profile

import (
	"fmt"

	"github.com/wbd2023/quill/internal/style"
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
	case style.ToolchainCheck:
		return validator.validateToolchainCheck(detail)

	case style.ProfileCheck:
		return validator.validateProfileCheck(detail)

	case style.FileCommand:
		return validator.validateFileCommand(detail)

	case style.TargetCommandTemplate:
		return validator.validateTargetCommand(detail)

	case style.TargetCheckTemplate:
		return validator.validateTargetCheck(detail)

	case style.RepositoryScan:
		return validator.validateRepositoryScan(detail)

	case style.ExternalCheck:
		return validator.validateExternalCheck(detail)

	default:
		return fmt.Errorf(
			"rule definition %q %s uses unknown execution template",
			validator.ruleID,
			validator.label,
		)
	}
}
