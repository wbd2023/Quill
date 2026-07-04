package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func absPath(path string) (resolved string, err error) {
	return filepath.Abs(path)
}

func resolveRepoRoot(path string) (repoRoot string, err error) {
	if path != "" {
		return absPath(path)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return findRepoRoot(workingDirectory)
}

func findRepoRoot(start string) (repoRoot string, err error) {
	directory, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		if isRepoRoot(directory) {
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

func isRepoRoot(path string) (root bool) {
	if fileExists(filepath.Join(path, "STYLE.md")) &&
		fileExists(filepath.Join(path, "style.toml")) {
		return true
	}

	return false
}

func fileExists(path string) (exists bool) {
	_, err := os.Stat(path)
	return err == nil
}
