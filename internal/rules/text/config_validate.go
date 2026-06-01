package text

import (
	"fmt"
	"slices"
	"strings"
)

// ValidateConfig validates text rule policy.
func ValidateConfig(config Config) (err error) {
	return validateSectionHeaders(config.SectionHeaders)
}

func validateSectionHeaders(headers SectionHeaderConfig) (err error) {
	if headers.LargeMinLines <= 0 {
		return fmt.Errorf("packs.text.section_headers.large_min_lines must be positive")
	}

	if headers.ShortMaxLines <= 0 {
		return fmt.Errorf("packs.text.section_headers.short_max_lines must be positive")
	}

	if headers.ShortMaxLines >= headers.LargeMinLines {
		field := "packs.text.section_headers.short_max_lines"
		return fmt.Errorf("%s must be less than large_min_lines", field)
	}

	if headers.MaxHeaderCount <= 0 {
		return fmt.Errorf("packs.text.section_headers.max_header_count must be positive")
	}

	return validateSectionHeaderNames(headers)
}

func validateSectionHeaderNames(headers SectionHeaderConfig) (err error) {
	if len(headers.GenericNames) == 0 {
		return fmt.Errorf("packs.text.section_headers.generic_names must not be empty")
	}

	names := slices.Concat(headers.GenericNames, headers.StructuralNames)
	seen := make(map[string]bool, len(names))
	for _, name := range names {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("packs.text.section_headers contains an empty header name")
		}

		if seen[name] {
			return fmt.Errorf(
				"packs.text.section_headers contains duplicate header name %q",
				name,
			)
		}

		seen[name] = true
	}

	return nil
}
