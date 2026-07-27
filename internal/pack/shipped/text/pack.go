package text

import (
	"github.com/wbd2023/quill/internal/checks/textpolicy"
	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/toolchain"
)

// PackID is the canonical identifier for this Pack.
const PackID = "text"

// Pack returns the Text Shipped Pack definition.
func Pack(tools []toolchain.Capability) (definition pack.Definition) {
	return pack.Definition{
		ID:       PackID,
		Name:     "Text",
		Tools:    append([]toolchain.Capability{}, tools...),
		FileSets: fileSets(),
		Config: pack.Config{
			Required: true,
			Validate: textpolicy.ValidatePackConfig,
		},
		Rules: rules(),
	}
}
