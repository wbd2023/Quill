package runner

import "ciphera/tools/internal/contract"

func CheckStatus(
	rule contract.Rule,
	err error,
	strictRecommendations bool,
) (status contract.CheckStatus) {
	switch {
	case err == nil:
		return contract.CheckStatusPass

	case IsBlocked(err):
		return contract.CheckStatusSkip

	case rule.Level == contract.LevelRecommendation && !strictRecommendations:
		return contract.CheckStatusWarn

	default:
		return contract.CheckStatusFail
	}
}
