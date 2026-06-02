package golang

import (
	"ciphera/tools/internal/pack"
	gopolicy "ciphera/tools/internal/rules/golang/policy"
	"ciphera/tools/internal/toolchain"
)

const (
	PackID = "go"

	ToolGo           = "go"
	ToolGoimports    = "goimports"
	ToolGolangciLint = "golangci-lint"
)

// Pack returns the Go Shipped Pack definition.
func Pack(tools []toolchain.Capability) (definition pack.Definition) {
	return pack.Definition{
		ID:    PackID,
		Name:  "Go",
		Tools: append([]toolchain.Capability{}, tools...),
		Config: pack.Config{
			Required: true,
			Validate: gopolicy.ValidatePackConfig,
		},
		Rules: rules(),
	}
}
