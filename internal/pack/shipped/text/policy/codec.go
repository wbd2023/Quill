package policy

import (
	"fmt"

	"github.com/wbd2023/quill/internal/profile"
)

// DecodeConfig decodes the Text Pack Policy subtree.
func DecodeConfig(policy profile.PackPolicy) (config Config, err error) {
	if policy == nil {
		return Config{}, fmt.Errorf("packs.text must be configured")
	}

	if err = rejectUnknownFields(policy, "packs.text", "section_headers"); err != nil {
		return Config{}, err
	}

	section, err := configSection(
		policy,
		"section_headers",
		"packs.text.section_headers",
	)
	if err != nil {
		return Config{}, err
	}

	config.SectionHeaders, err = decodeSectionHeaderConfig(section)
	if err != nil {
		return Config{}, err
	}

	return config, ValidateConfig(config)
}

// ValidatePackPolicy validates the raw Text Pack Policy subtree.
func ValidatePackPolicy(policy profile.PackPolicy) (err error) {
	_, err = DecodeConfig(policy)
	return err
}

// EncodeConfig encodes config as a raw Text Pack Policy subtree.
func EncodeConfig(config Config) (policy profile.PackPolicy) {
	return profile.PackPolicy{
		"section_headers": map[string]any{
			"large_min_lines":  config.SectionHeaders.LargeMinLines,
			"short_max_lines":  config.SectionHeaders.ShortMaxLines,
			"max_header_count": config.SectionHeaders.MaxHeaderCount,
			"generic_names":    cloneStrings(config.SectionHeaders.GenericNames),
			"structural_names": cloneStrings(config.SectionHeaders.StructuralNames),
		},
	}
}
