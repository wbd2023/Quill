package contract

const (
	EnforcementRequired       Enforcement = "required"
	EnforcementRecommendation Enforcement = "recommendation"
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

type Enforcement string

type CheckMode string

type Scope string

type CheckStatus string
