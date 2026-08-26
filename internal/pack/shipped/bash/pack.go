package bash

import (
	"github.com/wbd2023/quill/internal/pack"
)

// PackID is the canonical identifier for this Pack.
const PackID = "bash"

// Pack returns the Bash Shipped Pack definition. toolIDs reference the canonical Tool capabilities
// owned by the catalog by global ID.
func Pack(toolIDs ...string) (definition pack.Definition) {
	return pack.Definition{
		ID:       PackID,
		Name:     "Bash",
		ToolIDs:  append([]string{}, toolIDs...),
		FileSets: fileSets(),
		Rules:    rules(),
	}
}
