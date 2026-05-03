package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

func validateLanguage(
	repository policy.RepositoryConfig,
	language policy.LanguageConfig,
) (err error) {
	seen := make(map[string]bool, len(language.Backends))
	for _, backend := range language.Backends {
		if backend.Name == "" {
			return fmt.Errorf("language backend name must not be empty")
		}

		if seen[backend.Name] {
			return fmt.Errorf("duplicate language backend %q", backend.Name)
		}

		seen[backend.Name] = true

		if backend.Language == "" {
			return fmt.Errorf("language backend %q must define language", backend.Name)
		}

		if !repository.HasScope(backend.Scope) {
			return fmt.Errorf(
				"language backend %q references unknown scope %q",
				backend.Name,
				backend.Scope,
			)
		}
	}

	return nil
}
