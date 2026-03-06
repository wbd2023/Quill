package checker

import (
	"go/token"

	"stylecheck/internal/checker/collect"
)

/* -------------------------------------------- Types ------------------------------------------- */

// violation represents a single style rule violation.
type violation struct {
	position token.Position
	rule     string
	message  string
}

type analysisState struct {
	fileSet                *token.FileSet
	scannedGoFiles         []string
	violations             []violation
	interfaces             map[string]collect.InterfaceDecl
	mocks                  map[string][]collect.MethodDecl
	implementations        map[string][]collect.MethodDecl
	implementationBindings []collect.ImplementationBinding
}
