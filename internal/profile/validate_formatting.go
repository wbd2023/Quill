package profile

import (
	"fmt"
	"slices"

	"ciphera/tools/internal/policy"
)

func validateFormatting(formatting policy.FormattingConfig) (err error) {
	headers := formatting.SectionHeaders
	if headers.RequiredMinLines <= 0 {
		return fmt.Errorf("formatting.section_headers.required_min_lines must be positive")
	}

	if headers.ShortFileMaxLines <= 0 {
		return fmt.Errorf("formatting.section_headers.short_file_max_lines must be positive")
	}

	if headers.ShortFileMaxLines >= headers.RequiredMinLines {
		return fmt.Errorf(
			"formatting.section_headers.short_file_max_lines must be less than required_min_lines",
		)
	}

	if headers.OveruseThreshold <= 0 {
		return fmt.Errorf("formatting.section_headers.overuse_threshold must be positive")
	}

	return validateSectionHeaderNames(headers)
}

func validateSectionHeaderNames(headers policy.SectionHeaderConfig) (err error) {
	if len(headers.GenericNames) == 0 {
		return fmt.Errorf("formatting.section_headers.generic_names must not be empty")
	}

	names := slices.Concat(headers.GenericNames, headers.StructuralNames)
	seen := make(map[string]bool, len(names))
	for _, name := range names {
		if name == "" {
			return fmt.Errorf("formatting.section_headers contains an empty header name")
		}

		if seen[name] {
			return fmt.Errorf(
				"formatting.section_headers contains duplicate header name %q",
				name,
			)
		}

		seen[name] = true
	}

	return nil
}
