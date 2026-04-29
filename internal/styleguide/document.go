package styleguide

// Document is the parsed STYLE.md model used by coverage checks.
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
	Only   bool
	Reason string
}
