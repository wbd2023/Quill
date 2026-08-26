package report

import "github.com/wbd2023/quill/internal/style"

type fixJSON struct {
	Scope     string           `json:"scope"`
	Toolchain fixToolchainJSON `json:"toolchain"`
	Rules     []fixRuleJSON    `json:"rules"`
}

type fixToolchainJSON struct {
	AllValid bool             `json:"all_valid"`
	Statuses []toolStatusJSON `json:"statuses"`
}

type fixRuleJSON struct {
	RuleID         string            `json:"rule_id"`
	Name           string            `json:"name"`
	Group          style.RuleGroup   `json:"group"`
	Enforcement    style.Enforcement `json:"enforcement"`
	Scope          style.Scope       `json:"scope"`
	ExitCode       int               `json:"exit_code"`
	TimedOut       bool              `json:"timed_out"`
	Truncated      bool              `json:"truncated"`
	ExecutionError string            `json:"execution_error,omitempty"`
}

func newFixJSON(view FixView) (payload fixJSON) {
	return fixJSON{
		Scope: string(view.Scope),
		Toolchain: fixToolchainJSON{
			AllValid: view.AllValid,
			Statuses: toolStatusListJSON(view.Statuses),
		},
		Rules: fixRuleListJSON(view.Entries),
	}
}

func fixRuleListJSON(entries []FixEntry) (payload []fixRuleJSON) {
	payload = make([]fixRuleJSON, 0, len(entries))
	for _, entry := range entries {
		payload = append(payload, fixRuleJSON{
			RuleID:         entry.Rule.ID,
			Name:           entry.Rule.Name,
			Group:          entry.Rule.Group,
			Enforcement:    entry.Rule.Enforcement,
			Scope:          entry.Rule.Scope,
			ExitCode:       entry.Execution.ExitCode,
			TimedOut:       entry.Execution.TimedOut,
			Truncated:      entry.Execution.Truncated,
			ExecutionError: executionErrorMessage(entry.ExecutionError),
		})
	}

	return payload
}

func executionErrorMessage(err error) (message string) {
	if err == nil {
		return ""
	}

	return err.Error()
}
