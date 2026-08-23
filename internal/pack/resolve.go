package pack

import (
	"fmt"

	"github.com/wbd2023/quill/internal/profile"
)

/* --------------------------------------- Pack Resolution -------------------------------------- */

// ResolvePacks applies Pack-owned defaults and validates Pack-owned Profile Policy.
func ResolvePacks(
	config profile.Profile,
	packs []Definition,
) (resolved profile.Profile, err error) {
	resolved = config
	resolved.FileSets = resolveFileSets(config.FileSets, packs)

	if err = validatePackPolicies(resolved, packs); err != nil {
		return profile.Profile{}, err
	}

	return resolved, nil
}

func validatePackPolicies(config profile.Profile, packs []Definition) (err error) {
	active := indexPacks(packs)
	for packID := range config.PackPolicies {
		definition, found := active[packID]
		if !found {
			return fmt.Errorf("packs.%s policy is not active", packID)
		}

		if definition.Policy.Validate == nil {
			return fmt.Errorf("packs.%s policy is not supported", packID)
		}
	}

	for _, definition := range packs {
		packPolicy, found := config.PackPolicies.Lookup(definition.ID)
		if definition.Policy.Required && !found {
			return fmt.Errorf("packs.%s policy is required", definition.ID)
		}

		if !found || definition.Policy.Validate == nil {
			continue
		}

		if err = definition.Policy.Validate(packPolicy); err != nil {
			return err
		}
	}

	return nil
}

func indexPacks(packs []Definition) (indexed map[string]Definition) {
	indexed = make(map[string]Definition, len(packs))
	for _, definition := range packs {
		indexed[definition.ID] = definition
	}

	return indexed
}

/* ------------------------------------- File Set Resolution ------------------------------------ */

func resolveFileSets(
	configured profile.FileSets,
	packs []Definition,
) (fileSets profile.FileSets) {
	defaultCount := countDefaultFileSets(packs)
	if len(configured) == 0 && defaultCount == 0 {
		return nil
	}

	fileSets = make(profile.FileSets, 0, len(configured)+defaultCount)
	for _, definition := range packs {
		for _, fileSet := range definition.FileSets {
			fileSets = upsertFileSet(fileSets, fileSet.Clone())
		}
	}

	for _, fileSet := range configured {
		fileSets = upsertFileSet(fileSets, fileSet.Clone())
	}

	return fileSets
}

func countDefaultFileSets(packs []Definition) (count int) {
	for _, definition := range packs {
		count += len(definition.FileSets)
	}

	return count
}

func upsertFileSet(
	fileSets profile.FileSets,
	fileSet profile.FileSetConfig,
) (merged profile.FileSets) {
	for index := range fileSets {
		if fileSets[index].Name == fileSet.Name {
			fileSets[index] = fileSet
			return fileSets
		}
	}

	return append(fileSets, fileSet)
}
