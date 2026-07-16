package filewalk

import (
	"fmt"
	"strings"
)

// ValidateCollectorPolicy reports an error if config does not exclude standard directories or if
// the generated-file marker is empty.
func ValidateCollectorPolicy(config WalkConfig) (err error) {
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
		if isExcludedDirectory(config, directory) {
			continue
		}

		return fmt.Errorf("collector must exclude %q", directory)
	}

	if strings.TrimSpace(config.GeneratedMarker) == "" {
		return fmt.Errorf("collector generated-file marker must not be empty")
	}

	return nil
}
