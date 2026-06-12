package report

import "ciphera/tools/internal/style"

type CheckEntry struct {
	Rule   RuleSummary
	Status style.CheckStatus
	Result style.ExecutionResult
}

type RuleSummary struct {
	ID             string
	Name           string
	Group          style.RuleGroup
	Enforcement    style.Enforcement
	Scope          style.Scope
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

func (result CheckResult) Summary() (summary CheckSummary) {
	for _, entry := range result.Entries {
		switch entry.Status {
		case style.CheckStatusPass:
			summary.Passed++

		case style.CheckStatusWarn:
			summary.Warned++

		case style.CheckStatusFail:
			summary.Failed++

		case style.CheckStatusSkip:
			summary.Skipped++
		}
	}

	return summary
}
