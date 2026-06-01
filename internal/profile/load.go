package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/toml"
)

// DefaultFilename is the style profile filename loaded from repository roots.
const DefaultFilename = "style.toml"

// Load reads the default profile file from a repository root and validates it.
func Load(root string) (config policy.Config, err error) {
	path := filepath.Join(root, DefaultFilename)
	contents, err := os.ReadFile(path)
	if err != nil {
		return policy.Config{}, fmt.Errorf("read style profile %q: %w", path, err)
	}

	config, err = Parse(string(contents))
	if err != nil {
		return policy.Config{}, fmt.Errorf("load style profile %q: %w", path, err)
	}

	for _, marker := range config.Repository.RootMarkers {
		_, err = os.Stat(filepath.Join(root, marker))
		switch {
		case err == nil:
			continue

		case errors.Is(err, os.ErrNotExist):
			return policy.Config{}, fmt.Errorf(
				"repository root missing marker %q: %w",
				marker,
				err,
			)

		default:
			return policy.Config{}, fmt.Errorf(
				"repository root marker %q cannot be checked: %w",
				marker,
				err,
			)
		}
	}

	return config, nil
}

// Parse parses style profile TOML source and validates it.
func Parse(source string) (config policy.Config, err error) {
	config, err = toml.Decode(source)
	if err != nil {
		return policy.Config{}, err
	}

	if err = Validate(config); err != nil {
		return policy.Config{}, err
	}

	return config, nil
}
