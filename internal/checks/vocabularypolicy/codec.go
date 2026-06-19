package vocabularypolicy

import (
	"fmt"

	corepolicy "ciphera/tools/internal/policy"
)

// DecodeConfig decodes the Vocabulary Pack Policy subtree.
func DecodeConfig(pack corepolicy.PackConfig) (config Config, err error) {
	if pack == nil {
		return Config{}, fmt.Errorf("packs.vocabulary must be configured")
	}

	if err = rejectUnknownFields(pack, "packs.vocabulary", "go", "bash"); err != nil {
		return Config{}, err
	}

	goSection, err := configSection(pack, "go", "packs.vocabulary.go")
	if err != nil {
		return Config{}, err
	}

	bashSection, err := configSection(pack, "bash", "packs.vocabulary.bash")
	if err != nil {
		return Config{}, err
	}

	config.Go, err = decodeGoConfig(goSection)
	if err != nil {
		return Config{}, err
	}

	config.Bash, err = decodeBashConfig(bashSection)
	if err != nil {
		return Config{}, err
	}

	return config, ValidateConfig(config)
}

// ValidatePackConfig validates the raw Vocabulary Pack Policy subtree.
func ValidatePackConfig(pack corepolicy.PackConfig) (err error) {
	_, err = DecodeConfig(pack)
	return err
}

// EncodeConfig encodes config as a raw Vocabulary Pack Policy subtree.
func EncodeConfig(config Config) (pack corepolicy.PackConfig) {
	return corepolicy.PackConfig{
		"go": map[string]any{
			"type_suffixes":       encodeStringListMap(config.Go.TypeSuffixes),
			"identifier_suffixes": encodeStringListMap(config.Go.IdentifierSuffixes),
		},
		"bash": map[string]any{
			"variable_names": encodeStringListMap(config.Bash.VariableNames),
		},
	}
}
