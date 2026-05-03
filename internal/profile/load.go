package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/policy"

	"github.com/BurntSushi/toml"
)

const defaultFilename = "style.toml"

// Load reads the default profile file from repoRoot and validates it.
func Load(repoRoot string) (config policy.Config, err error) {
	config, err = loadFile(filepath.Join(repoRoot, defaultFilename))
	if err != nil {
		return policy.Config{}, err
	}

	if err = Validate(config); err != nil {
		return policy.Config{}, err
	}

	if err = validateRepositoryRoot(repoRoot, config.Repository); err != nil {
		return policy.Config{}, fmt.Errorf(
			"repository root does not satisfy profile markers: %w",
			err,
		)
	}

	return config, nil
}

func validateRepositoryRoot(repoRoot string, repository policy.RepositoryConfig) (err error) {
	for _, marker := range repository.RootMarkers {
		if marker == "" {
			return fmt.Errorf("repository root marker must not be empty")
		}

		if _, statErr := os.Stat(filepath.Join(repoRoot, marker)); statErr != nil {
			if errors.Is(statErr, os.ErrNotExist) {
				return fmt.Errorf("repository root missing marker %q: %w", marker, statErr)
			}

			return fmt.Errorf("repository root marker %q cannot be checked: %w", marker, statErr)
		}
	}

	return nil
}

func loadFile(path string) (config policy.Config, err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return policy.Config{}, err
	}

	return parse(string(contents))
}

func parse(contents string) (config policy.Config, err error) {
	var schema schemaConfig
	metadata, err := toml.Decode(contents, &schema)
	if err != nil {
		return policy.Config{}, err
	}

	undecoded := metadata.Undecoded()
	if len(undecoded) > 0 {
		return policy.Config{}, fmt.Errorf("unknown style.toml key %q", formatTOMLKey(undecoded[0]))
	}

	return configFromSchema(schema), nil
}

func formatTOMLKey(key toml.Key) (text string) {
	parts := make([]string, 0, len(key))
	for _, part := range key {
		parts = append(parts, part)
	}

	return strings.Join(parts, ".")
}
