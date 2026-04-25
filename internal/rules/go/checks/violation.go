package checks

import "go/token"

// Violation represents a single style rule violation.
type Violation struct {
	Position token.Position
	Rule     string
	Message  string
}
