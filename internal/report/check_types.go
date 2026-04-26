package report

import "ciphera/tools/internal/contract"

/* ---------------------------------------- Check Entries --------------------------------------- */

type CheckEntry struct {
	Rule   RuleSummary
	Status contract.CheckStatus
	Result contract.ExecutionResult
}

type RuleSummary struct {
	ID             string
	Name           string
	Group          contract.RuleGroup
	Level          contract.Level
	Scope          contract.Scope
	RequirementIDs []string
}

type CheckResult struct {
	Entries []CheckEntry
}

type CheckSummary struct {
	Passed  int
	Warned  int
	Failed  int
	Skipped int
}

type CheckGroup struct {
	Group   contract.RuleGroup
	Entries []CheckEntry
}

type CheckView struct {
	Result  CheckResult
	Summary CheckSummary
	Groups  []CheckGroup
}

/* ------------------------------------------ JSON DTOs ----------------------------------------- */

type checkJSON struct {
	Result  checkResultJSON  `json:"result"`
	Summary CheckSummary     `json:"summary"`
	Groups  []checkGroupJSON `json:"groups"`
}

type checkResultJSON struct {
	Entries []checkEntryJSON `json:"entries"`
}

type checkGroupJSON struct {
	Group   contract.RuleGroup `json:"group"`
	Entries []checkEntryJSON   `json:"entries"`
}

type checkEntryJSON struct {
	RuleID       string               `json:"rule_id"`
	Name         string               `json:"name"`
	Group        contract.RuleGroup   `json:"group"`
	Level        contract.Level       `json:"level"`
	Scope        contract.Scope       `json:"scope"`
	Status       contract.CheckStatus `json:"status"`
	Requirements []string             `json:"requirements"`
	Diagnostics  []diagnosticJSON     `json:"diagnostics"`
	Output       string               `json:"output,omitempty"`
	Command      *commandResultJSON   `json:"command,omitempty"`
}

type diagnosticJSON struct {
	Code    string `json:"code"`
	File    string `json:"file,omitempty"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	Message string `json:"message"`
}

type commandResultJSON struct {
	ExitCode  int  `json:"exit_code"`
	TimedOut  bool `json:"timed_out"`
	Truncated bool `json:"truncated"`
}

/* --------------------------------------- Check Summaries -------------------------------------- */

func (result CheckResult) Summary() (summary CheckSummary) {
	for _, entry := range result.Entries {
		switch entry.Status {
		case contract.CheckStatusPass:
			summary.Passed++

		case contract.CheckStatusWarn:
			summary.Warned++

		case contract.CheckStatusFail:
			summary.Failed++

		case contract.CheckStatusSkip:
			summary.Skipped++
		}
	}

	return summary
}
