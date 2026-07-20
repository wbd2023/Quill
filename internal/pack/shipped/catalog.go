package shipped

import (
	"github.com/wbd2023/Quill/internal/pack"
	"github.com/wbd2023/Quill/internal/pack/shipped/bash"
	"github.com/wbd2023/Quill/internal/pack/shipped/golang"
	"github.com/wbd2023/Quill/internal/pack/shipped/markdown"
	"github.com/wbd2023/Quill/internal/pack/shipped/project"
	"github.com/wbd2023/Quill/internal/pack/shipped/security"
	"github.com/wbd2023/Quill/internal/pack/shipped/text"
	"github.com/wbd2023/Quill/internal/pack/shipped/tool"
	"github.com/wbd2023/Quill/internal/pack/shipped/vocabulary"
)

// DefaultCatalog returns the Shipped Pack catalogue.
func DefaultCatalog() (catalog pack.Catalog) {
	return pack.NewCatalog(
		project.Pack(tool.BuildAll()),
		text.Pack(tool.Select(tool.Misspell)),
		markdown.Pack(tool.Select(tool.Markdownlint)),
		bash.Pack(tool.Select(
			tool.Shellcheck,
			tool.Shfmt,
		)),
		golang.Pack(tool.Select(
			tool.Go,
			tool.Goimports,
			tool.GolangciLint,
		)),
		security.Pack(),
		vocabulary.Pack(),
	)
}

// DefaultRegistry builds a registry from the Shipped Pack catalogue.
func DefaultRegistry(enabled []string) (registry pack.Registry, err error) {
	return DefaultCatalog().Registry(enabled)
}
