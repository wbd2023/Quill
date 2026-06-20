package bash

import (
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/toolchain"
)

// PackID is the canonical identifier for this Pack.
const PackID = "bash"

// Pack returns the Bash Shipped Pack definition.
func Pack(tools []toolchain.Capability) (definition pack.Definition) {
	return pack.Definition{
		ID:       PackID,
		Name:     "Bash",
		Tools:    append([]toolchain.Capability{}, tools...),
		FileSets: fileSets(),
		Rules:    rules(),
	}
}
