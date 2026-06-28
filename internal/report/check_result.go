package report

import "ciphera/tools/internal/style"

// CheckEntry is check entry.
type CheckEntry struct {
	Rule   RuleSummary
	Status style.CheckStatus
	Result style.ExecutionResult
}

// RuleSummary is rule summary.
type RuleSummary struct {
	ID             string
	Name           string
	Group          style.RuleGroup
	Enforcement    style.Enforcement
	Scope          style.Scope
	RequirementIDs []string
}

// CheckResult is check result.
type CheckResult struct {
	Entries []CheckEntry
}

// CheckSummary is check summary.
type CheckSummary struct {
	Passed  int
	Warned  int
	Failed  int
	Skipped int
	Errored int
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

		case style.CheckStatusError:
			summary.Errored++
		}
	}

	return summary
}
