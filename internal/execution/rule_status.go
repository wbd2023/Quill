package execution

import "github.com/wbd2023/Quill/internal/style"

// CheckStatus classifies the outcome of running one rule. Findings are data carried by result: a
// non-empty result means the rule found violations; a non-nil error means the rule could not run
// (operational failure). The two are distinct - a parse error is not a style failure.
func CheckStatus(
	rule style.Rule,
	result style.ExecutionResult,
	err error,
	strictRecommendations bool,
) (status style.CheckStatus) {
	switch {
	case IsBlocked(err):
		return style.CheckStatusSkip

	case err != nil:
		return style.CheckStatusError

	case result.Empty():
		return style.CheckStatusPass

	case rule.Enforcement == style.EnforcementRecommendation && !strictRecommendations:
		return style.CheckStatusWarn

	default:
		return style.CheckStatusFail
	}
}
