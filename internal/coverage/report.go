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

type Status string

type Mode string

type Requirement struct {
	ID      string
	Section string
	Text    string
	Mode    Mode
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
