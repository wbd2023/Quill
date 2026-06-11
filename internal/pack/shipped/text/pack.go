package text

import (
	textpolicy "ciphera/tools/internal/checks/text/policy"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/toolchain"
)

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
