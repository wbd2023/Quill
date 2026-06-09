package project

import (
	"fmt"

	projectrules "ciphera/tools/internal/checks/project"
	"ciphera/tools/internal/runner"
)

func decodeProjectConfig(
	context runner.Context,
	packID string,
) (config projectrules.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return projectrules.Config{}, errMissingPackConfig(packID)
	}

	return projectrules.DecodeConfig(pack)
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
