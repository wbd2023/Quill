package text

import (
	"github.com/wbd2023/quill/internal/pack"
	textpolicy "github.com/wbd2023/quill/internal/pack/shipped/text/policy"
)

// PackID is the canonical identifier for this Pack.
const PackID = "text"

// Pack returns the Text Shipped Pack definition. toolIDs reference the canonical Tool capabilities
// owned by the catalog by global ID.
func Pack(toolIDs ...string) (definition pack.Definition) {
	return pack.Definition{
		ID:       PackID,
		Name:     "Text",
		ToolIDs:  append([]string{}, toolIDs...),
		FileSets: fileSets(),
		Policy: pack.Policy{
			Required: true,
			Validate: textpolicy.Validate,
		},
		Rules: rules(),
	}
}
