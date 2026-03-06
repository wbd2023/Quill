package collect

import "go/token"

type MethodDecl struct {
	Name     string
	Position token.Position
}

type InterfaceDecl struct {
	Name     string
	Methods  []MethodDecl
	Position token.Position
}

type ImplementationBinding struct {
	InterfaceName      string
	ImplementationName string
	ImplementationKey  string
	Position           token.Position
}
