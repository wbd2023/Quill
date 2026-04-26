package coverage

import "ciphera/tools/internal/styleguide"

const (
	StatusAutomated  Status = "automated"
	StatusPartial    Status = "partial"
	StatusReviewOnly Status = "review_only"
	StatusManual     Status = "manual"
)

type Status string

type Requirement struct {
	ID      string
	Section string
	Text    string
	Mode    styleguide.VerificationMode
	Reason  string
	RuleIDs []string
}

type Section struct {
	Section             string
	Title               string
	Status              Status
	RequirementCount    int
	AutomatedCount      int
	ReviewOnlyCount     int
	ManualDeferredCount int
}

type Report struct {
	Requirements []Requirement
	Sections     []Section
}
