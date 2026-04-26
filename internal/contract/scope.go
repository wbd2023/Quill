package contract

const (
	LevelRequired       Level = "required"
	LevelRecommendation Level = "recommendation"
)

const (
	CheckModeRequired CheckMode = "required"
	CheckModeAll      CheckMode = "all"
)

const (
	CheckStatusPass CheckStatus = "pass"
	CheckStatusWarn CheckStatus = "warn"
	CheckStatusFail CheckStatus = "fail"
	CheckStatusSkip CheckStatus = "skip"
)

type Level string

type CheckMode string

type Scope string

type CheckStatus string
