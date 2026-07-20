package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/checks/projectpolicy"
	"github.com/wbd2023/Quill/internal/execution"
)

func decodeProjectConfig(
	context execution.RunContext,
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
