package runner

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/report"
)

func CheckStatus(
	rule contract.Rule,
	err error,
	strictRecommendations bool,
) (status report.CheckStatus) {
	switch {
	case err == nil:
		return report.CheckStatusPass

	case IsBlocked(err):
		return report.CheckStatusSkip

	case rule.Level == contract.LevelRecommendation && !strictRecommendations:
		return report.CheckStatusWarn

	default:
		return report.CheckStatusFail
	}
}
