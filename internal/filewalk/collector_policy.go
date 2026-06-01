package filewalk

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/policy"
)

func ValidateCollectorPolicy(repository policy.RepositoryConfig) (err error) {
	requiredDirectories := []string{
		".cache",
		".git",
		".toolchain",
		"bin",
		"testdata",
		"third_party",
		"vendor",
	}

	for _, directory := range requiredDirectories {
		if isExcludedDirectory(repository, directory) {
			continue
		}

		return fmt.Errorf("collector must exclude %q", directory)
	}

	if strings.TrimSpace(repository.GeneratedMarker) == "" {
		return fmt.Errorf("collector generated-file marker must not be empty")
	}

	return nil
}
