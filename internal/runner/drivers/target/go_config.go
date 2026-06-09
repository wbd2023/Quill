package target

import (
	"fmt"

	gopolicy "ciphera/tools/internal/checks/golang/policy"
	"ciphera/tools/internal/runner"
)

func decodeGoConfig(
	context runner.Context,
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
