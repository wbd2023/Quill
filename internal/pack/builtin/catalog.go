package builtin

import (
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/builtin/bash"
	"ciphera/tools/internal/pack/builtin/golang"
	"ciphera/tools/internal/pack/builtin/markdown"
	"ciphera/tools/internal/pack/builtin/project"
	"ciphera/tools/internal/pack/builtin/security"
	"ciphera/tools/internal/pack/builtin/text"
	"ciphera/tools/internal/pack/builtin/vocabulary"
)

// DefaultCatalog returns the Shipped Pack catalogue.
func DefaultCatalog() (catalog pack.Catalog) {
	return pack.NewCatalog(
		project.Pack(coreTools()),
		text.Pack(selectTools(text.ToolMisspell)),
		markdown.Pack(selectTools(markdown.ToolMarkdownlint)),
		bash.Pack(selectTools(
			bash.ToolShellcheck,
			bash.ToolShfmt,
		)),
		golang.Pack(selectTools(
			golang.ToolGo,
			golang.ToolGoimports,
			golang.ToolGolangciLint,
		)),
		security.Pack(),
		vocabulary.Pack(),
	)
}

// DefaultRegistry builds a registry from the Shipped Pack catalogue.
func DefaultRegistry(enabled []string) (registry pack.Registry, err error) {
	return DefaultCatalog().Registry(enabled)
}
