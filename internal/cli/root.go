package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func resolveRepositoryRoot(path string) (repositoryRoot string, err error) {
	if path != "" {
		return filepath.Abs(path)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return findRepositoryRoot(workingDirectory)
}

func findRepositoryRoot(start string) (repositoryRoot string, err error) {
	directory, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		if isRepositoryRoot(directory) {
			return directory, nil
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}

		directory = parent
	}

	return "", fmt.Errorf("could not locate repository root from %q", start)
}

func isRepositoryRoot(path string) (root bool) {
	if pathExists(filepath.Join(path, "STYLE.md")) &&
		pathExists(filepath.Join(path, "quill.toml")) {
		return true
	}

	return false
}

func pathExists(path string) (exists bool) {
	_, err := os.Stat(path)
	return err == nil
}
