package styleguide

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	VerificationAutomated      VerificationMode = "automated"
	VerificationReviewOnly     VerificationMode = "review_only"
	VerificationManualDeferred VerificationMode = "manual_deferred"
)

const (
	CoverageAutomated  CoverageStatus = "automated"
	CoveragePartial    CoverageStatus = "partial"
	CoverageReviewOnly CoverageStatus = "review_only"
	CoverageManual     CoverageStatus = "manual"
)

/* -------------------------------------------- Types ------------------------------------------- */

type VerificationMode string

type CoverageStatus string

type RequirementMetadata struct {
	ID     string
	Mode   VerificationMode
	Reason string
}

type Requirement struct {
	ID      string
	Section string
	Text    string
	Mode    VerificationMode
	Reason  string
	RuleIDs []string
}

type SectionCoverage struct {
	Section             string
	Title               string
	Status              CoverageStatus
	RequirementCount    int
	AutomatedCount      int
	ReviewOnlyCount     int
	ManualDeferredCount int
}

type CoverageReport struct {
	Requirements []Requirement
	Sections     []SectionCoverage
}

/* ------------------------------------- Coverage Reporting ------------------------------------- */

func Coverage(repoRoot string) (report CoverageReport, err error) {
	headings, err := readHeadings(repoRoot)
	if err != nil {
		return CoverageReport{}, err
	}

	requirements, err := Requirements(repoRoot)
	if err != nil {
		return CoverageReport{}, err
	}

	return CoverageReport{
		Requirements: requirements,
		Sections:     buildSectionCoverage(headings, requirements),
	}, nil
}

func buildSectionCoverage(
	headings []documentHeading,
	requirements []Requirement,
) (sections []SectionCoverage) {
	requirementsBySection := make(map[string][]Requirement)
	for _, requirement := range requirements {
		requirementsBySection[requirement.Section] = append(
			requirementsBySection[requirement.Section],
			requirement,
		)
	}

	sections = make([]SectionCoverage, 0, len(headings))
	for _, heading := range headings {
		entry := SectionCoverage{
			Section: heading.Section,
			Title:   heading.Title,
		}

		for _, requirement := range requirementsBySection[heading.Section] {
			entry.RequirementCount++

			switch requirement.Mode {
			case VerificationAutomated:
				entry.AutomatedCount++
			case VerificationReviewOnly:
				entry.ReviewOnlyCount++
			case VerificationManualDeferred:
				entry.ManualDeferredCount++
			}
		}

		entry.Status = deriveCoverageStatus(entry)
		sections = append(sections, entry)
	}

	return sections
}

func deriveCoverageStatus(entry SectionCoverage) (status CoverageStatus) {
	switch {
	case entry.RequirementCount == 0:
		return CoverageManual

	case entry.AutomatedCount == entry.RequirementCount:
		return CoverageAutomated

	case entry.AutomatedCount > 0:
		return CoveragePartial

	case entry.ReviewOnlyCount == entry.RequirementCount:
		return CoverageReviewOnly

	default:
		return CoverageManual
	}
}
