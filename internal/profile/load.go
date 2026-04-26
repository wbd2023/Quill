package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/policy"

	"github.com/BurntSushi/toml"
)

func Load(repoRoot string) (config policy.Config, err error) {
	config, err = loadProfileFile(filepath.Join(repoRoot, "style.toml"))
	if err != nil {
		return policy.Config{}, err
	}

	if err = Validate(config); err != nil {
		return policy.Config{}, err
	}

	if err = config.Repository.ValidateRoot(repoRoot); err != nil {
		return policy.Config{}, fmt.Errorf(
			"repository root does not satisfy profile markers: %w",
			err,
		)
	}

	return config, nil
}

func loadProfileFile(path string) (config policy.Config, err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return policy.Config{}, err
	}

	return parseProfile(string(contents))
}

func parseProfile(contents string) (config policy.Config, err error) {
	var schema schemaConfig
	metadata, err := toml.Decode(contents, &schema)
	if err != nil {
		return policy.Config{}, err
	}

	undecoded := metadata.Undecoded()
	if len(undecoded) > 0 {
		return policy.Config{}, fmt.Errorf("unknown style.toml key %q", undecodedKey(undecoded[0]))
	}

	return policyFromSchema(schema), nil
}

func undecodedKey(key toml.Key) (text string) {
	parts := make([]string, 0, len(key))
	for _, part := range key {
		parts = append(parts, part)
	}

	return strings.Join(parts, ".")
}
