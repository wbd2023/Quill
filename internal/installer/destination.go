package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func prepareDestinationDirectory(
	root string,
	destination string,
) (path string, directory string, err error) {
	root, err = filepath.Abs(root)
	if err != nil {
		return "", "", fmt.Errorf("resolve installation root %q: %w", root, err)
	}
	path, err = filepath.Abs(destination)
	if err != nil {
		return "", "", fmt.Errorf("resolve installation destination %q: %w", destination, err)
	}

	relative, err := filepath.Rel(root, path)
	if err != nil {
		return "", "", fmt.Errorf("resolve destination beneath root: %w", err)
	}

	if relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return "", "", fmt.Errorf("installation destination %q escapes root %q", path, root)
	}

	directory = root
	parent := filepath.Dir(relative)
	if parent == "." {
		return path, directory, nil
	}

	for _, component := range strings.Split(parent, string(os.PathSeparator)) {
		if component == "" || component == "." {
			continue
		}

		directory = filepath.Join(directory, component)
		info, statErr := os.Lstat(directory)
		if os.IsNotExist(statErr) {
			mkdirErr := os.Mkdir(directory, standardPermissions)
			if mkdirErr != nil && !os.IsExist(mkdirErr) {
				return "", "", fmt.Errorf(
					"create installation directory %q: %w",
					directory,
					mkdirErr,
				)
			}
			info, statErr = os.Lstat(directory)
		}
		if statErr != nil {
			return "", "", fmt.Errorf("inspect installation directory %q: %w", directory, statErr)
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return "", "", fmt.Errorf("installation directory %q is a symlink", directory)
		}

		if !info.IsDir() {
			return "", "", fmt.Errorf("installation path %q is not a directory", directory)
		}
	}

	return path, directory, nil
}

func prepareExecutableDestination(
	root string,
	destination string,
) (path string, directory string, exists bool, err error) {
	path, directory, err = prepareDestinationDirectory(root, destination)
	if err != nil {
		return "", "", false, err
	}

	info, statErr := os.Lstat(path)
	if os.IsNotExist(statErr) {
		return path, directory, false, nil
	}

	if statErr != nil {
		return "", "", false, fmt.Errorf("inspect destination %q: %w", path, statErr)
	}

	if !info.Mode().IsRegular() {
		return "", "", false, fmt.Errorf("refuse to use non-regular destination %q", path)
	}

	return path, directory, true, nil
}
