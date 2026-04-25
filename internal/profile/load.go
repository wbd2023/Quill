package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

/* ------------------------------------------- Loading ------------------------------------------ */

func Load(repoRoot string) (policy Profile, err error) {
	policy, err = loadProfileFile(filepath.Join(repoRoot, "style.toml"))
	if err != nil {
		return Profile{}, err
	}

	if err = policy.Validate(); err != nil {
		return Profile{}, err
	}

	if err = policy.Repository.ValidateRoot(repoRoot); err != nil {
		return Profile{}, fmt.Errorf("repository root does not satisfy profile markers: %w", err)
	}

	return policy, nil
}

func loadProfileFile(path string) (policy Profile, err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, err
	}

	return parseProfile(string(contents))
}

func parseProfile(contents string) (policy Profile, err error) {
	policy = Profile{
		Paths: PathClassSet{Classes: make(map[string][]string)},
	}

	metadata, err := toml.Decode(contents, &policy)
	if err != nil {
		return Profile{}, err
	}

	undecoded := metadata.Undecoded()
	if len(undecoded) > 0 {
		return Profile{}, fmt.Errorf("unknown style.toml key %q", undecodedKey(undecoded[0]))
	}

	return policy, nil
}

func undecodedKey(key toml.Key) (text string) {
	parts := make([]string, 0, len(key))
	for _, part := range key {
		parts = append(parts, part)
	}

	return strings.Join(parts, ".")
}
