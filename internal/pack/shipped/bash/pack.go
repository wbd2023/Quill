package bash

import (
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/toolchain"
)

// PackID is pack i d.
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
