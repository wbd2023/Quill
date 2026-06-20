package coverage

const (
	StatusAutomated  Status = "automated"
	StatusPartial    Status = "partial"
	StatusReviewOnly Status = "review_only"
	StatusManual     Status = "manual"
)

const (
	ModeAutomated      Mode = "automated"
	ModeReviewOnly     Mode = "review_only"
	ModeManualDeferred Mode = "manual_deferred"
)

// Status is status.
type Status string

// Mode is mode.
type Mode string

// Requirement is requirement.
type Requirement struct {
	ID      string
	Section string
	Text    string
	Mode    Mode
	Reason  string
	RuleIDs []string
}

// Section is section.
type Section struct {
	Section             string
	Title               string
	Status              Status
	RequirementCount    int
	AutomatedCount      int
	ReviewOnlyCount     int
	ManualDeferredCount int
}

// Report is report.
type Report struct {
	Requirements []Requirement
	Sections     []Section
}
