package builtin

import "ciphera/tools/internal/pack"

// DefaultCatalog returns the built-in Ciphera Pack catalog.
func DefaultCatalog() (catalog pack.Catalog) {
	return pack.NewCatalog(
		projectPack(),
		textPack(),
		markdownPack(),
		bashPack(),
		goPack(),
		securityPack(),
		vocabularyPack(),
	)
}

// DefaultRegistry builds a registry from the built-in Ciphera Pack catalog.
func DefaultRegistry(enabled []string) (registry pack.Registry, err error) {
	return DefaultCatalog().Registry(enabled)
}
