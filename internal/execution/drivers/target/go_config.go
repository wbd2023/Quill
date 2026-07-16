package target

import (
	"fmt"

	"ciphera/tools/internal/checks/gopolicy"
	"ciphera/tools/internal/execution"
)

func decodeGoConfig(
	context execution.Context,
	packID string,
) (config gopolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return gopolicy.Config{}, errMissingPackConfig(packID)
	}

	return gopolicy.DecodeConfig(pack)
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
