package order

import "go/token"

type methodDeclaration struct {
	Name     string
	Position token.Position
}

type interfaceDeclaration struct {
	Name     string
	Methods  []methodDeclaration
	Position token.Position
}

type implementationBinding struct {
	InterfaceName      string
	ImplementationName string
	ImplementationKey  string
	Position           token.Position
}
