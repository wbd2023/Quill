package effective

import (
	"fmt"

	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/policy"
)

/* --------------------------------------- Pack Resolution -------------------------------------- */

// ResolvePacks applies pack-owned defaults and validates pack-owned profile config.
func ResolvePacks(
	config policy.Config,
	packs []pack.Definition,
) (resolved policy.Config, err error) {
	resolved = config
	resolved.FileSets = resolveFileSets(config.FileSets, packs)

	if err = validatePackConfigs(resolved, packs); err != nil {
		return policy.Config{}, err
	}

	return resolved, nil
}

func validatePackConfigs(config policy.Config, packs []pack.Definition) (err error) {
	active := indexPacks(packs)
	for packID := range config.PackConfigs {
		pack, found := active[packID]
		if !found {
			return fmt.Errorf("packs.%s config is not active", packID)
		}

		if pack.Config.Validate == nil {
			return fmt.Errorf("packs.%s config is not supported", packID)
		}
	}

	for _, pack := range packs {
		packConfig, found := config.PackConfigs.Lookup(pack.ID)
		if pack.Config.Required && !found {
			return fmt.Errorf("packs.%s must be configured", pack.ID)
		}

		if !found || pack.Config.Validate == nil {
			continue
		}

		if err = pack.Config.Validate(packConfig); err != nil {
			return err
		}
	}

	return nil
}

func indexPacks(packs []pack.Definition) (indexed map[string]pack.Definition) {
	indexed = make(map[string]pack.Definition, len(packs))
	for _, pack := range packs {
		indexed[pack.ID] = pack
	}

	return indexed
}

/* ------------------------------------- File Set Resolution ------------------------------------ */

func resolveFileSets(
	configured policy.FileSets,
	packs []pack.Definition,
) (fileSets policy.FileSets) {
	defaultCount := countDefaultFileSets(packs)
	if len(configured) == 0 && defaultCount == 0 {
		return nil
	}

	fileSets = make(policy.FileSets, 0, len(configured)+defaultCount)
	for _, pack := range packs {
		for _, fileSet := range pack.FileSets {
			fileSets = upsertFileSet(fileSets, fileSet.Clone())
		}
	}

	for _, fileSet := range configured {
		fileSets = upsertFileSet(fileSets, fileSet.Clone())
	}

	return fileSets
}

func countDefaultFileSets(packs []pack.Definition) (count int) {
	for _, pack := range packs {
		count += len(pack.FileSets)
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
