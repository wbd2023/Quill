package profile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// DefaultFilename is the Quill Profile filename loaded from repository roots.
const DefaultFilename = "quill.toml"

// Load reads the default profile through a repository-confined filesystem and
// validates it.
func Load(root string) (config Profile, err error) {
	path := filepath.Join(root, DefaultFilename)
	repository, err := os.OpenRoot(root)
	if err != nil {
		return Profile{}, fmt.Errorf(
			"open repository root %q: %w",
			root,
			err,
		)
	}

	contents, readErr := fs.ReadFile(repository.FS(), DefaultFilename)
	closeErr := repository.Close()
	if readErr != nil {
		return Profile{}, fmt.Errorf(
			"read style profile %q: %w",
			path,
			readErr,
		)
	}

	if closeErr != nil {
		return Profile{}, fmt.Errorf(
			"close repository root %q: %w",
			root,
			closeErr,
		)
	}

	config, err = Parse(string(contents))
	if err != nil {
		return Profile{}, fmt.Errorf("load style profile %q: %w", path, err)
	}

	if err = validateRepositoryPaths(config, root); err != nil {
		return Profile{}, fmt.Errorf("load style profile %q: %w", path, err)
	}

	for _, marker := range config.Repository.RootMarkers {
		_, err = os.Stat(filepath.Join(root, marker))
		switch {
		case err == nil:
			continue

		case errors.Is(err, os.ErrNotExist):
			return Profile{}, fmt.Errorf(
				"repository root missing marker %q: %w",
				marker,
				err,
			)

		default:
			return Profile{}, fmt.Errorf(
				"repository root marker %q cannot be checked: %w",
				marker,
				err,
			)
		}
	}

	return config, nil
}

// Parse parses style profile TOML source and validates it.
func Parse(source string) (config Profile, err error) {
	config, err = decodeTOML(source)
	if err != nil {
		return Profile{}, err
	}

	if err = Validate(config); err != nil {
		return Profile{}, err
	}

	return config, nil
}
