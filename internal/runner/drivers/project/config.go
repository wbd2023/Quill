package project

import (
	"fmt"

	"ciphera/tools/internal/pack/builtin"
	projectrules "ciphera/tools/internal/rules/project"
	"ciphera/tools/internal/runner"
)

func decodeProjectConfig(context runner.Context) (config projectrules.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(builtin.PackProject)
	if !found {
		return projectrules.Config{}, errMissingPackConfig(builtin.PackProject)
	}

	return projectrules.DecodeConfig(pack)
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
