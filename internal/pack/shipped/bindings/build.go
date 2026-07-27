// Package bindings is the explicit aggregate composition point for shipped runtime behaviour.
//
// It does not own any Pack-specific registration details. Each shipped Pack owns a child bindings
// package that wires its own execution identities; this package composes them in one place so the
// engine can construct a complete drivers.Bindings value without global registration or init side
// effects.
package bindings

import (
	"github.com/wbd2023/quill/internal/execution/drivers"
	bashbindings "github.com/wbd2023/quill/internal/pack/shipped/bash/bindings"
	golangbindings "github.com/wbd2023/quill/internal/pack/shipped/golang/bindings"
	markdownbindings "github.com/wbd2023/quill/internal/pack/shipped/markdown/bindings"
	projectbindings "github.com/wbd2023/quill/internal/pack/shipped/project/bindings"
	securitybindings "github.com/wbd2023/quill/internal/pack/shipped/security/bindings"
	textbindings "github.com/wbd2023/quill/internal/pack/shipped/text/bindings"
	vocabularybindings "github.com/wbd2023/quill/internal/pack/shipped/vocabulary/bindings"
)

// Build composes every shipped Pack's runtime bindings into a single drivers.Bindings value for
// driver construction. Each Pack's registrations live in its own child bindings package; this
// function is composition only and owns no Pack-specific details.
func Build() (bindings drivers.Bindings) {
	bindings = drivers.NewBindings()
	bashbindings.Register(&bindings)
	golangbindings.Register(&bindings)
	markdownbindings.Register(&bindings)
	projectbindings.Register(&bindings)
	securitybindings.Register(&bindings)
	textbindings.Register(&bindings)
	vocabularybindings.Register(&bindings)
	return bindings
}
