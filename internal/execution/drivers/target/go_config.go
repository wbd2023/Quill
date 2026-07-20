package target

import (
	"fmt"
	"strings"

	"github.com/wbd2023/Quill/internal/checks/gopolicy"
	"github.com/wbd2023/Quill/internal/execution"
)

func decodeGoConfig(
	context execution.RunContext,
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

func joinGoLocalImportPrefixes(prefixes []string) (prefix string) {
	return strings.Join(prefixes, ",")
}
