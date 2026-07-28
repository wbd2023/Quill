package golang

import (
	"github.com/wbd2023/quill/internal/pack"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
)

// PackID is the canonical identifier for this Pack.
const PackID = "go"

// Pack returns the Go Shipped Pack definition. toolIDs reference the canonical Tool capabilities
// owned by the catalogue by global ID.
func Pack(toolIDs ...string) (definition pack.Definition) {
	return pack.Definition{
		ID:      PackID,
		Name:    "Go",
		ToolIDs: append([]string{}, toolIDs...),
		Config: pack.Config{
			Required: true,
			Validate: gopolicy.ValidatePackConfig,
		},
		Rules: rules(),
	}
}
