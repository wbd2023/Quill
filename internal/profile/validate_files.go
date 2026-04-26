package profile

import (
	"fmt"

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

	if headers.OveruseCount <= 0 {
		return fmt.Errorf("formatting.section_headers.overuse_header_count must be positive")
	}

	if len(headers.GenericNames) == 0 {
		return fmt.Errorf("formatting.section_headers.generic_names must not be empty")
	}

	seen := make(map[string]bool, len(headers.GenericNames)+len(headers.StructuralNames))
	for _, names := range [][]string{headers.GenericNames, headers.StructuralNames} {
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
	}

	return nil
}

func validateFileSets(
	repository policy.RepositoryConfig,
	fileSets []policy.FileSetConfig,
) (err error) {
	seenFileSets := make(map[string]bool, len(fileSets))
	for _, fileSet := range fileSets {
		if fileSet.Name == "" {
			return fmt.Errorf("file set name must not be empty")
		}

		if seenFileSets[fileSet.Name] {
			return fmt.Errorf("duplicate file set %q", fileSet.Name)
		}

		seenFileSets[fileSet.Name] = true
		for scope := range fileSet.Files {
			if !repository.ScopeExists(scope) {
				return fmt.Errorf("file set %q references unknown scope %q", fileSet.Name, scope)
			}
		}

		for scope := range fileSet.Prefixes {
			if !repository.ScopeExists(scope) {
				return fmt.Errorf("file set %q references unknown scope %q", fileSet.Name, scope)
			}
		}
	}

	return nil
}
