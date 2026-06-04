package scan

import (
	"fmt"

	"ciphera/tools/internal/pack/builtin"
	gopolicy "ciphera/tools/internal/rules/golang/policy"
	"ciphera/tools/internal/rules/text"
	"ciphera/tools/internal/rules/vocabulary"
	"ciphera/tools/internal/runner"
)

func decodeGoPackConfig(context runner.Context) (config gopolicy.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(builtin.PackGo)
	if !found {
		return gopolicy.Config{}, errMissingPackConfig(builtin.PackGo)
	}

	return gopolicy.DecodeConfig(pack)
}

func decodeTextPackConfig(context runner.Context) (config text.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(builtin.PackText)
	if !found {
		return text.Config{}, errMissingPackConfig(builtin.PackText)
	}

	return text.DecodeConfig(pack)
}

func decodeVocabularyPackConfig(
	context runner.Context,
) (config vocabulary.Config, err error) {
	pack, found := context.Profile.PackConfigs.Lookup(builtin.PackVocabulary)
	if !found {
		return vocabulary.Config{}, errMissingPackConfig(builtin.PackVocabulary)
	}

	return vocabulary.DecodeConfig(pack)
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
