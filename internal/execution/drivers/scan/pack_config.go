package scan

import (
	"fmt"

	"ciphera/tools/internal/checks/gopolicy"
	"ciphera/tools/internal/checks/textpolicy"
	"ciphera/tools/internal/checks/vocabularypolicy"
	"ciphera/tools/internal/execution"
)

func decodeGoPackConfig(
	context execution.Context,
	packID string,
) (config gopolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return gopolicy.Config{}, errMissingPackConfig(packID)
	}

	return gopolicy.DecodeConfig(pack)
}

func decodeTextPackConfig(
	context execution.Context,
	packID string,
) (config textpolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(packID)
	if !found {
		return textpolicy.Config{}, errMissingPackConfig(packID)
	}

	return textpolicy.DecodeConfig(pack)
}

func decodeVocabularyPackConfig(
	context execution.Context,
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
