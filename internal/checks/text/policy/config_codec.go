package policy

import (
	"fmt"

	corepolicy "ciphera/tools/internal/policy"
)

// DecodeConfig decodes the Text Pack Policy subtree.
func DecodeConfig(pack corepolicy.PackConfig) (config Config, err error) {
	if pack == nil {
		return Config{}, fmt.Errorf("packs.text must be configured")
	}

	if err = rejectUnknownFields(pack, "packs.text", "section_headers"); err != nil {
		return Config{}, err
	}

	section, err := configSection(
		pack,
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

// ValidatePackConfig validates the raw Text Pack Policy subtree.
func ValidatePackConfig(pack corepolicy.PackConfig) (err error) {
	_, err = DecodeConfig(pack)
	return err
}

// EncodeConfig encodes config as a raw Text Pack Policy subtree.
func EncodeConfig(config Config) (pack corepolicy.PackConfig) {
	return corepolicy.PackConfig{
		"section_headers": map[string]any{
			"large_min_lines":  config.SectionHeaders.LargeMinLines,
			"short_max_lines":  config.SectionHeaders.ShortMaxLines,
			"max_header_count": config.SectionHeaders.MaxHeaderCount,
			"generic_names":    cloneStrings(config.SectionHeaders.GenericNames),
			"structural_names": cloneStrings(config.SectionHeaders.StructuralNames),
		},
	}
}
