package profile

import (
	"fmt"

	"ciphera/tools/internal/checks/projectpolicy"
	"ciphera/tools/internal/runner"
)

func decodeProjectConfig(
	context runner.Context,
	packID string,
) (config projectpolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return projectpolicy.Config{}, errMissingPackConfig(packID)
	}

	return projectpolicy.DecodeConfig(pack)
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
