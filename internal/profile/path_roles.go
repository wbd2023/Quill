package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/policy"
)

func validatePathRoles(paths policy.PathRoles) (err error) {
	for role, patterns := range paths {
		if isBlank(role) {
			return fmt.Errorf("path_roles contains an empty path role")
		}

		if len(patterns) == 0 {
			return fmt.Errorf("path_roles.%s must not be empty", role)
		}

		seen := make(map[string]bool, len(patterns))
		for _, pattern := range patterns {
			if isBlank(pattern) {
				return fmt.Errorf("path_roles.%s contains an empty pattern", role)
			}

			if seen[pattern] {
				return fmt.Errorf("path_roles.%s duplicates pattern %q", role, pattern)
			}

			seen[pattern] = true
		}
	}

	return nil
}
