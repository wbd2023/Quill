package styleguide

const (
	VerificationAutomated      VerificationMode = "automated"
	VerificationReviewOnly     VerificationMode = "review_only"
	VerificationManualDeferred VerificationMode = "manual_deferred"
)

type VerificationMode string

type RequirementMetadata struct {
	ID     string
	Mode   VerificationMode
	Reason string
}
