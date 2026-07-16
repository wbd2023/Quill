package pack

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

/* --------------------------------------- Pack Resolution -------------------------------------- */

// ResolvePacks applies pack-owned defaults and validates pack-owned profile config.
func ResolvePacks(
	config policy.Config,
	packs []Definition,
) (resolved policy.Config, err error) {
	resolved = config
	resolved.FileSets = resolveFileSets(config.FileSets, packs)

	if err = validatePackConfigs(resolved, packs); err != nil {
		return policy.Config{}, err
	}

	return resolved, nil
}

func validatePackConfigs(config policy.Config, packs []Definition) (err error) {
	active := indexPacks(packs)
	for packID := range config.PackConfigs {
		definition, found := active[packID]
		if !found {
			return fmt.Errorf("packs.%s config is not active", packID)
		}

		if definition.Config.Validate == nil {
			return fmt.Errorf("packs.%s config is not supported", packID)
		}
	}

	for _, definition := range packs {
		packConfig, found := config.PackConfigs.Lookup(definition.ID)
		if definition.Config.Required && !found {
			return fmt.Errorf("packs.%s must be configured", definition.ID)
		}

		if !found || definition.Config.Validate == nil {
			continue
		}

		if err = definition.Config.Validate(packConfig); err != nil {
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
	configured policy.FileSets,
	packs []Definition,
) (fileSets policy.FileSets) {
	defaultCount := countDefaultFileSets(packs)
	if len(configured) == 0 && defaultCount == 0 {
		return nil
	}

	fileSets = make(policy.FileSets, 0, len(configured)+defaultCount)
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
	fileSets policy.FileSets,
	fileSet policy.FileSetConfig,
) (merged policy.FileSets) {
	for index := range fileSets {
		if fileSets[index].Name == fileSet.Name {
			fileSets[index] = fileSet
			return fileSets
		}
	}

	return append(fileSets, fileSet)
}
