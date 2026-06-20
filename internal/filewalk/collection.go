package filewalk

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

/* ----------------------------------------- Collection ----------------------------------------- */

// CollectFiles collect files.
func CollectFiles(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope style.Scope,
	extensions ...string,
) (paths []string, err error) {
	roots := repository.ResolveScopeRoots(repoRoot, scope)
	return collectFilesInRoots(roots, repository, func(path string) bool {
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

// CollectFilesInRoots collect files in roots.
func CollectFilesInRoots(
	repository policy.RepositoryConfig,
	roots []string,
	extensions ...string,
) (paths []string, err error) {
	return collectFilesInRoots(roots, repository, func(path string) bool {
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

// CollectFilesInScopes collect files in scopes.
func CollectFilesInScopes(
	repoRoot string,
	repository policy.RepositoryConfig,
	scopes []style.Scope,
	extensions ...string,
) (paths []string, err error) {
	return CollectFilesInRoots(
		repository,
		collectScopeRoots(repoRoot, repository, scopes),
		extensions...,
	)
}

// CollectAllFiles collect all files.
func CollectAllFiles(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope style.Scope,
) (paths []string, err error) {
	roots := repository.ResolveScopeRoots(repoRoot, scope)
	return collectFilesInRoots(roots, repository, func(path string) bool {
		info, statErr := os.Stat(path)
		if statErr != nil {
			return false
		}

		return info.Mode().IsRegular()
	})
}

func collectFilesInRoots(
	roots []string,
	repository policy.RepositoryConfig,
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

				if entry.IsDir() && isExcludedDirectory(repository, entry.Name()) {
					return filepath.SkipDir
				}

				if entry.IsDir() {
					return nil
				}

				if isGeneratedFile(path, repository) {
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
