package pack

import (
	"fmt"
	"strings"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/toolchain"
)

/* ------------------------------------------ Catalogue ----------------------------------------- */

// Catalog stores the canonical Tool capabilities and the Packs available to a style checker build.
// The catalogue owns each Tool capability exactly once by global ID; Packs reference Tools by ID
// instead of carrying copies, so a Tool is never defined twice and a Pack can never silently
// override another Pack's Tool declaration.
type Catalog struct {
	tools []toolchain.Capability
	packs []Definition
}

// NewCatalog returns a catalogue owning defensive deep copies of the canonical tools and the
// supplied packs. The capabilities are the single source of Tool definitions; packs reference them
// by ID.
func NewCatalog(tools []toolchain.Capability, packs ...Definition) (catalog Catalog) {
	return Catalog{
		tools: CloneCapabilities(tools),
		packs: CloneDefinitions(packs),
	}
}

// Tools returns defensive deep copies of the canonical Tool capabilities, including mutable
// installer data.
func (catalog Catalog) Tools() (capabilities []toolchain.Capability) {
	return CloneCapabilities(catalog.tools)
}

// Packs returns defensive copies of the packs available in the catalog.
func (catalog Catalog) Packs() (packs []Definition) {
	return CloneDefinitions(catalog.packs)
}

/* ------------------------------------------ Registry ------------------------------------------ */

// Registry builds a rule registry by resolving each selected Pack's Tool references against the
// canonical Tool capabilities. Duplicate Tool declarations and unknown references are rejected
// before any loss, and every rule is stamped with its declaring Pack's ID for provenance.
func (catalog Catalog) Registry(enabled []string) (registry Registry, err error) {
	if err = validateCatalog(catalog); err != nil {
		return Registry{}, err
	}

	packs := catalog.Packs()
	if len(enabled) > 0 {
		if packs, err = selectPacks(packs, enabled); err != nil {
			return Registry{}, err
		}
	}

	if registry, err = buildRegistry(catalog.Tools(), packs); err != nil {
		return Registry{}, err
	}

	if err = validateRegistry(registry); err != nil {
		return Registry{}, err
	}

	return registry, nil
}

/* ----------------------------------------- Validation ----------------------------------------- */

func validateCatalog(catalog Catalog) (err error) {
	seenPacks := make(map[string]bool, len(catalog.packs))
	for _, pack := range catalog.packs {
		if strings.TrimSpace(pack.ID) == "" {
			return fmt.Errorf("catalog contains an empty pack id")
		}
		if pack.ID == profile.EnabledPacksKey {
			return fmt.Errorf("catalog contains reserved pack id %q", pack.ID)
		}
		if strings.TrimSpace(pack.Name) == "" {
			return fmt.Errorf("catalog contains pack %q with an empty name", pack.ID)
		}

		if seenPacks[pack.ID] {
			return fmt.Errorf("catalog contains duplicate pack id %q", pack.ID)
		}

		seenPacks[pack.ID] = true
	}

	seenRules := make(map[string]string)
	for _, pack := range catalog.packs {
		for _, rule := range pack.Rules {
			if strings.TrimSpace(rule.ID) == "" {
				return fmt.Errorf("catalog contains an empty rule id in pack %q", pack.ID)
			}
			if owner := seenRules[rule.ID]; owner != "" {
				return fmt.Errorf(
					"catalog contains duplicate rule id %q (declared by packs %q and %q)",
					rule.ID, owner, pack.ID,
				)
			}
			seenRules[rule.ID] = pack.ID
		}
	}

	seenTools := make(map[string]bool, len(catalog.tools))
	for _, tool := range catalog.tools {
		if tool.ID == "" {
			return fmt.Errorf("catalog contains a tool with an empty id")
		}

		if tool.Name == "" {
			return fmt.Errorf("catalog contains a tool %q with an empty name", tool.ID)
		}

		if seenTools[tool.ID] {
			return fmt.Errorf("catalog contains duplicate tool id %q", tool.ID)
		}

		seenTools[tool.ID] = true
	}

	return nil
}
