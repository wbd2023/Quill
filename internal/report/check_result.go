package report

import "github.com/wbd2023/quill/internal/style"

// CheckEntry is one rendered Rule outcome.
type CheckEntry struct {
	Rule           RuleSummary
	Status         style.CheckStatus
	Result         style.ExecutionResult
	ExecutionError error
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

// CheckSummary is the aggregate outcome of a check.
type CheckSummary struct {
	Passed  int
	Warned  int
	Failed  int
	Blocked int
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

		case style.CheckStatusBlocked:
			summary.Blocked++
		case style.CheckStatusSkip:
			summary.Skipped++

		case style.CheckStatusError:
			summary.Errored++
		}
	}

	return summary
}

// hasCommandMetadata reports whether the execution result carries command output or metadata.
// It preserves the historical facade semantics where captured command output counts as a command,
// even when the exit code is zero and neither timeout nor truncation occurred.
func hasCommandMetadata(result style.ExecutionResult) (hasMetadata bool) {
	return result.Output != "" || result.HasCommand()
}
