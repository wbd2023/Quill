package filewalk

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

/* ----------------------------------------- Collection ----------------------------------------- */

// CollectFiles collects files under roots that match any of extensions. If no extensions are
// given, every regular file is collected.
func CollectFiles(
	roots []string,
	config WalkConfig,
	extensions ...string,
) (paths []string, err error) {
	return collectFilesInRoots(roots, config, func(path string) bool {
		if len(extensions) == 0 {
			return true
		}

		for _, extension := range extensions {
			if strings.HasSuffix(path, extension) {
				return true
			}
		}

		return false
	})
}

// CollectFilesInRoots collects files in roots that match any of extensions. If no extensions are
// given, every regular file is collected.
func CollectFilesInRoots(
	config WalkConfig,
	roots []string,
	extensions ...string,
) (paths []string, err error) {
	return collectFilesInRoots(roots, config, func(path string) bool {
		if len(extensions) == 0 {
			info, statErr := os.Stat(path)
			if statErr != nil {
				return false
			}

			return info.Mode().IsRegular()
		}

		for _, extension := range extensions {
			if strings.HasSuffix(path, extension) {
				return true
			}
		}

		return false
	})
}

// CollectAllFiles collects all regular files under roots.
func CollectAllFiles(roots []string, config WalkConfig) (paths []string, err error) {
	return collectFilesInRoots(roots, config, func(path string) bool {
		info, statErr := os.Stat(path)
		if statErr != nil {
			return false
		}

		return info.Mode().IsRegular()
	})
}

func collectFilesInRoots(
	roots []string,
	config WalkConfig,
	include func(path string) bool,
) (paths []string, err error) {
	for _, root := range roots {
		if _, statErr := os.Stat(root); statErr != nil {
			continue
		}

		walkErr := filepath.WalkDir(
			root,
			func(path string, entry fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}

				if entry.IsDir() && isExcludedDirectory(config, entry.Name()) {
					return filepath.SkipDir
				}

				if entry.IsDir() {
					return nil
				}

				if isGeneratedFile(path, config.GeneratedMarker) {
					return nil
				}

				if !include(path) {
					return nil
				}

				paths = append(paths, filepath.Clean(path))
				return nil
			},
		)
		if walkErr != nil {
			return nil, walkErr
		}
	}

	sort.Strings(paths)
	return dedupePaths(paths), nil
}
