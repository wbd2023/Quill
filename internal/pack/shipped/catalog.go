package shipped

import (
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped/bash"
	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/markdown"
	"ciphera/tools/internal/pack/shipped/project"
	"ciphera/tools/internal/pack/shipped/security"
	"ciphera/tools/internal/pack/shipped/text"
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/pack/shipped/vocabulary"
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
