package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

func validateFileSets(
	repository policy.RepositoryConfig,
	fileSets []policy.FileSetConfig,
) (err error) {
	seen := make(map[string]bool, len(fileSets))
	for _, fileSet := range fileSets {
		if fileSet.Name == "" {
			return fmt.Errorf("file set name must not be empty")
		}

		if seen[fileSet.Name] {
			return fmt.Errorf("duplicate file set %q", fileSet.Name)
		}

		seen[fileSet.Name] = true
		for scope := range fileSet.ExplicitFiles {
			if !repository.HasScope(scope) {
				return fmt.Errorf("file set %q references unknown scope %q", fileSet.Name, scope)
			}
		}

		for scope := range fileSet.PathPrefixes {
			if !repository.HasScope(scope) {
				return fmt.Errorf("file set %q references unknown scope %q", fileSet.Name, scope)
			}
		}
	}

	return nil
}
