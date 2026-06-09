package runner

import "ciphera/tools/internal/style"

func CheckStatus(
	rule style.Rule,
	err error,
	strictRecommendations bool,
) (status style.CheckStatus) {
	switch {
	case err == nil:
		return style.CheckStatusPass

	case IsBlocked(err):
		return style.CheckStatusSkip

	case rule.Enforcement == style.EnforcementRecommendation && !strictRecommendations:
		return style.CheckStatusWarn

	default:
		return style.CheckStatusFail
	}
}
