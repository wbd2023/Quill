package drivers

import (
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/project"
	"ciphera/tools/internal/runner"
)

func decodeProjectConfig(context runner.Context) (config project.Config, err error) {
	pack, found := context.Policy.PackConfigs.Lookup(builtin.PackProject)
	if !found {
		return project.Config{}, errMissingPackConfig(builtin.PackProject)
	}

	return project.DecodeConfig(pack)
}
