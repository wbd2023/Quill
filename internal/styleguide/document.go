package styleguide

// Document is the parsed STYLE.md model for requirement binding and coverage.
type Document struct {
	Headings     []Heading
	Requirements []Requirement
}

// Heading is a numbered STYLE.md section heading.
type Heading struct {
	Section string
	Title   string
}

// Requirement is a documented STYLE.md requirement.
type Requirement struct {
	ID      string
	Section string
	Text    string
	Review  Review
}

// Review describes review-only metadata for a STYLE.md requirement.
type Review struct {
	IsReviewOnly bool
	Reason       string
}
