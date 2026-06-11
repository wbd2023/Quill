package policy

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
			"forbidden_type_suffixes":       cloneStrings(config.Go.ForbiddenTypeSuffixes),
			"preferred_type_suffix":         config.Go.PreferredTypeSuffix,
			"forbidden_identifier_suffixes": cloneStrings(config.Go.ForbiddenIdentifierSuffixes),
			"preferred_identifier_suffix":   config.Go.PreferredIdentifierSuffix,
		},
		"bash": map[string]any{
			"forbidden_variable_names": cloneStrings(config.Bash.ForbiddenVariableNames),
			"preferred_variable_name":  config.Bash.PreferredVariableName,
		},
	}
}
