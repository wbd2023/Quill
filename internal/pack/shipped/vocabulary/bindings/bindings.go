// Package bindings owns the Vocabulary Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Vocabulary execution identities (the vocabulary repository
// scanner) to concrete generic drivers. The parent vocabulary package stays independent of the
// driver facade.
package bindings

import (
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/vocabulary"
)

// Register wires every Vocabulary execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(
		vocabulary.ScannerVocabulary,
		drivers.CheckVocabulary(vocabulary.PackID),
	)
}
