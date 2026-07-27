package golang

import (
	"github.com/wbd2023/quill/internal/checks/gopolicy"
	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/toolchain"
)

// PackID is the canonical identifier for this Pack.
const PackID = "go"

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
