package shipped

import (
	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/pack/shipped/bash"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
	"github.com/wbd2023/quill/internal/pack/shipped/markdown"
	"github.com/wbd2023/quill/internal/pack/shipped/project"
	"github.com/wbd2023/quill/internal/pack/shipped/security"
	"github.com/wbd2023/quill/internal/pack/shipped/text"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/pack/shipped/vocabulary"
)

// DefaultCatalog returns the Shipped Pack catalogue. The catalogue owns every canonical Tool
// capability exactly once; each Pack references its Tools by global ID.
func DefaultCatalog() (catalog pack.Catalog) {
	return pack.NewCatalog(
		tool.BuildAll(),
		project.Pack(
			tool.Go,
			tool.Goimports,
			tool.Misspell,
			tool.GolangciLint,
			tool.Shfmt,
			tool.Shellcheck,
			tool.Markdownlint,
		),
		text.Pack(tool.Misspell),
		markdown.Pack(tool.Markdownlint),
		bash.Pack(
			tool.Shellcheck,
			tool.Shfmt,
		),
		golang.Pack(
			tool.Go,
			tool.Goimports,
			tool.GolangciLint,
		),
		security.Pack(),
		vocabulary.Pack(),
	)
}

// DefaultRegistry builds a registry from the Shipped Pack catalogue.
func DefaultRegistry(enabled []string) (registry pack.Registry, err error) {
	return DefaultCatalog().Registry(enabled)
}

// ComposeCatalog returns the Shipped catalogue augmented with external Pack definitions. External
// Packs declare no canonical Tools - their Rules reference self-describing ExternalCheck Jobs - so
// composition does not create a second selection, binding, or execution pipeline.
func ComposeCatalog(external []pack.Definition) (catalog pack.Catalog) {
	base := DefaultCatalog()
	return pack.NewCatalog(base.Tools(), append(base.Packs(), external...)...)
}
