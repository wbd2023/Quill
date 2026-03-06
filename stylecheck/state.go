package main

import "go/token"

/* -------------------------------------------- Types ------------------------------------------- */

// violation represents a single style rule violation.
type violation struct {
	position token.Position
	rule     string
	message  string
}

type methodDecl struct {
	name     string
	position token.Position
}

type interfaceDecl struct {
	name     string
	methods  []methodDecl
	position token.Position
}

type implementationBinding struct {
	interfaceName      string
	implementationName string
	implementationKey  string
	position           token.Position
}

type analysisState struct {
	fileSet                *token.FileSet
	scannedGoFiles         []string
	violations             []violation
	interfaces             map[string]interfaceDecl
	mocks                  map[string][]methodDecl
	implementations        map[string][]methodDecl
	implementationBindings []implementationBinding
}
