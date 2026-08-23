package policy

import (
	"fmt"

	"github.com/wbd2023/quill/internal/profile"
)

// DecodeConfig decodes the Vocabulary Pack Policy subtree.
func DecodeConfig(policy profile.PackPolicy) (config Config, err error) {
	if policy == nil {
		return Config{}, fmt.Errorf("packs.vocabulary must be configured")
	}

	if err = rejectUnknownFields(policy, "packs.vocabulary", "go", "bash"); err != nil {
		return Config{}, err
	}

	goSection, err := configSection(policy, "go", "packs.vocabulary.go")
	if err != nil {
		return Config{}, err
	}

	bashSection, err := configSection(policy, "bash", "packs.vocabulary.bash")
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

// ValidatePackPolicy validates the raw Vocabulary Pack Policy subtree.
func ValidatePackPolicy(policy profile.PackPolicy) (err error) {
	_, err = DecodeConfig(policy)
	return err
}

// EncodeConfig encodes config as a raw Vocabulary Pack Policy subtree.
func EncodeConfig(config Config) (policy profile.PackPolicy) {
	return profile.PackPolicy{
		"go": map[string]any{
			"type_suffixes":       encodeStringListMap(config.Go.TypeSuffixes),
			"identifier_suffixes": encodeStringListMap(config.Go.IdentifierSuffixes),
		},
		"bash": map[string]any{
			"variable_names": encodeStringListMap(config.Bash.VariableNames),
		},
	}
}
