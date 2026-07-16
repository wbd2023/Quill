package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

func validatePackConfigs(
	enabledPacks []string,
	configs policy.PackConfigs,
) (err error) {
	if len(configs) == 0 {
		return nil
	}

	enabled := make(map[string]bool, len(enabledPacks))
	for _, packID := range enabledPacks {
		enabled[packID] = true
	}

	for packID, config := range configs {
		if isBlank(packID) {
			return fmt.Errorf("packs contains an empty pack id")
		}

		if !enabled[packID] {
			return fmt.Errorf("packs.%s config is not enabled in packs.enabled", packID)
		}

		if len(config) == 0 {
			return fmt.Errorf("packs.%s must not be empty", packID)
		}
	}

	return nil
}
