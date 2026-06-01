package pack

import "fmt"

// Catalog stores the packs available to a style checker build.
type Catalog struct {
	packs []Definition
}

// NewCatalog returns a catalog containing the supplied packs.
func NewCatalog(packs ...Definition) (catalog Catalog) {
	catalog.packs = CloneDefinitions(packs)
	return catalog
}

// Packs returns the packs available in the catalog.
func (catalog Catalog) Packs() (packs []Definition) {
	return CloneDefinitions(catalog.packs)
}

// Registry builds a rule registry from the catalog.
func (catalog Catalog) Registry(enabled []string) (registry Registry, err error) {
	packs := catalog.Packs()
	if err = validateCatalog(packs); err != nil {
		return Registry{}, err
	}

	if len(enabled) > 0 {
		packs, err = selectPacks(packs, enabled)
		if err != nil {
			return Registry{}, err
		}
	}

	registry = buildRegistry(packs)
	if err = validateRegistry(registry); err != nil {
		return Registry{}, err
	}

	return registry, nil
}

func validateCatalog(packs []Definition) (err error) {
	seen := make(map[string]bool, len(packs))
	for _, pack := range packs {
		if pack.ID == "" {
			return fmt.Errorf("catalog contains an empty pack id")
		}

		if seen[pack.ID] {
			return fmt.Errorf("catalog contains duplicate pack id %q", pack.ID)
		}

		seen[pack.ID] = true
	}

	return nil
}
