package scan

import (
	"fmt"

	gopolicy "ciphera/tools/internal/checks/golang/policy"
	"ciphera/tools/internal/checks/text"
	"ciphera/tools/internal/checks/vocabulary"
	"ciphera/tools/internal/runner"
)

func decodeGoPackConfig(context runner.Context, packID string) (config gopolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return gopolicy.Config{}, errMissingPackConfig(packID)
	}

	return gopolicy.DecodeConfig(pack)
}

func decodeTextPackConfig(context runner.Context, packID string) (config text.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return text.Config{}, errMissingPackConfig(packID)
	}

	return text.DecodeConfig(pack)
}

func decodeVocabularyPackConfig(
	context runner.Context,
	packID string,
) (config vocabulary.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return vocabulary.Config{}, errMissingPackConfig(packID)
	}

	return vocabulary.DecodeConfig(pack)
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
