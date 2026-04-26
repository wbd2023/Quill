package filewalk

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* ----------------------------------------- Collection ----------------------------------------- */

func CollectFiles(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope contract.Scope,
	extensions ...string,
) (paths []string, err error) {
	roots := scopeRoots(repoRoot, repository, scope)
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

func CollectFilesInScopes(
	repoRoot string,
	repository policy.RepositoryConfig,
	scopes []contract.Scope,
	extensions ...string,
) (paths []string, err error) {
	return CollectFilesInRoots(
		repository,
		repository.ScanRootsForScopes(repoRoot, scopes),
		extensions...,
	)
}

func CollectAllFiles(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope contract.Scope,
) (paths []string, err error) {
	roots := scopeRoots(repoRoot, repository, scope)
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
