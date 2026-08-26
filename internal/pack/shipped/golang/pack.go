package golang

import (
	"github.com/wbd2023/quill/internal/pack"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
	"github.com/wbd2023/quill/internal/profile"
)

// PackID is the canonical identifier for this Pack.
const PackID = "go"

// Pack returns the Go Shipped Pack definition. toolIDs reference the canonical Tool capabilities
// owned by the catalog by global ID.
func Pack(toolIDs ...string) (definition pack.Definition) {
	return pack.Definition{
		ID:      PackID,
		Name:    "Go",
		ToolIDs: append([]string{}, toolIDs...),
		Policy: pack.Policy{
			Required: true,
			Validate: func(policy profile.PackPolicy) error {
				_, err := gopolicy.DecodeConfig(policy)
				return err
			},
		},
		Rules: rules(),
	}
}
