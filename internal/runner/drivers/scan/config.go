package scan

import (
	"fmt"

	gopolicy "ciphera/tools/internal/checks/golang/policy"
	textpolicy "ciphera/tools/internal/checks/text/policy"
	vocabularypolicy "ciphera/tools/internal/checks/vocabulary/policy"
	"ciphera/tools/internal/runner"
)

func decodeGoPackConfig(context runner.Context, packID string) (config gopolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return gopolicy.Config{}, errMissingPackConfig(packID)
	}

	return gopolicy.DecodeConfig(pack)
}

func decodeTextPackConfig(
	context runner.Context,
	packID string,
) (config textpolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return textpolicy.Config{}, errMissingPackConfig(packID)
	}

	return textpolicy.DecodeConfig(pack)
}

func decodeVocabularyPackConfig(
	context runner.Context,
	packID string,
) (config vocabularypolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return vocabularypolicy.Config{}, errMissingPackConfig(packID)
	}

	return vocabularypolicy.DecodeConfig(pack)
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
