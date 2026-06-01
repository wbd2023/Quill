package drivers

import (
	"ciphera/tools/internal/pack/builtin"
	gopolicy "ciphera/tools/internal/rules/golang/policy"
	"ciphera/tools/internal/runner"
)

func decodeGoConfig(context runner.Context) (config gopolicy.Config, err error) {
	pack, found := context.Policy.PackConfigs.Lookup(builtin.PackGo)
	if !found {
		return gopolicy.Config{}, errMissingPackConfig(builtin.PackGo)
	}

	return gopolicy.DecodeConfig(pack)
}
