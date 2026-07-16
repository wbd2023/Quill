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
	// CheckModeRequired means only MUST-level rules run.
	CheckModeRequired CheckMode = "required"
	// CheckModeAll means both MUST-level and recommendation rules run.
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
	// CheckStatusError means the rule could not run (operational failure: parse error, missing
	// config, IO failure). Distinct from Fail, which means the rule ran and found violations.
	CheckStatusError CheckStatus = "error"
)

// Enforcement represents a rule's enforcement level: required or recommendation.
type Enforcement string

// CheckMode represents which rules a check pass runs.
type CheckMode string

// Scope represents a repository area, for example app, tools, or all.
type Scope string

// CheckStatus represents the per-rule outcome of a check run.
type CheckStatus string
