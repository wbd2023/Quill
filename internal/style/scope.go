package style

// Enforcement levels control whether a rule failure blocks CI or is reported only.
const (
	// EnforcementRequired means a rule failure fails the build.
	EnforcementRequired Enforcement = "required"
	// EnforcementRecommendation means a rule failure is reported but does not fail.
	EnforcementRecommendation Enforcement = "recommendation"
)

// CheckMode selects which rules run: only MUST-level rules or all rules.
const (
	// CheckModeRequired runs only MUST-level rules.
	CheckModeRequired CheckMode = "required"
	// CheckModeAll runs both MUST-level and recommendation rules.
	CheckModeAll CheckMode = "all"
)

// CheckStatus is the per-rule outcome of a check run.
const (
	// CheckStatusPass means the rule found no violations.
	CheckStatusPass CheckStatus = "pass"
	// CheckStatusWarn means the rule found recommendation-level findings.
	CheckStatusWarn CheckStatus = "warn"
	// CheckStatusFail means the rule found required-level violations.
	CheckStatusFail CheckStatus = "fail"
	// CheckStatusSkip means the rule was skipped due to a missing tool or blocked dependency.
	CheckStatusSkip CheckStatus = "skip"
)

// Enforcement is a rule's enforcement level: required or recommendation.
type Enforcement string

// CheckMode selects which rules run in a check pass.
type CheckMode string

// Scope names a repository area, for example app, tools, or all.
type Scope string

// CheckStatus is the per-rule outcome of a check run.
type CheckStatus string
