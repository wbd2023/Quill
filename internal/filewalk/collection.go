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
// given, every regular file is collected. Symlinks and other non-regular leaves (FIFO, socket,
// device) are excluded before any content is read so they never reach a generated/binary probe.
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

// CollectAllFiles collects all non-binary regular files under roots. Symlinks and other
// non-regular leaves are excluded before any content is read.
func CollectAllFiles(roots []string, config WalkConfig) (paths []string, err error) {
	return collectFilesInRoots(roots, config, func(path string) bool {
		return !IsBinaryFile(path)
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

				if entry.Name() == ".git" {
					if entry.IsDir() {
						return filepath.SkipDir
					}

					return nil
				}

				if entry.IsDir() && isExcludedDirectory(config, entry.Name()) {
					return filepath.SkipDir
				}

				if entry.IsDir() {
					return nil
				}

				// Reject symlink and other non-regular leaves (FIFO, socket,
				// device) using the walk entry's lstat type. This runs before
				// any generated/binary probe reads the leaf, so a link or
				// special file is excluded rather than opened or followed.
				if !entry.Type().IsRegular() {
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

/* ------------------------------------- Leaf Classification ------------------------------------ */

// IsRegularLeaf reports whether path is a regular, non-symlink leaf using lstat, so symlinks,
// FIFOs, sockets, and devices all report false. Collection already enforces this invariant
// through the walk entry; this helper covers explicit candidate paths that bypass the walk,
// letting callers fail closed before reading a candidate or handing it off to a driver.
func IsRegularLeaf(path string) (regular bool) {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}
