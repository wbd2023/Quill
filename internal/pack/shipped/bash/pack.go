package bash

import (
	"github.com/wbd2023/Quill/internal/pack"
	"github.com/wbd2023/Quill/internal/toolchain"
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
