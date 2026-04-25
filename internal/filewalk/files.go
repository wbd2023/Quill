package filewalk

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const generatedHeaderLineLimit = 12

var generatedCommentPrefixes = []string{"//", "#", ";", "--"}

/* ----------------------------------------- Validation ----------------------------------------- */

func ValidateCollectorPolicy(repository profile.RepositoryConfig) (err error) {
	requiredDirectories := []string{
		".cache",
		".git",
		".toolchain",
		"bin",
		"testdata",
		"third_party",
		"vendor",
	}

	for _, directory := range requiredDirectories {
		if isExcludedDirectory(repository, directory) {
			continue
		}

		return fmt.Errorf("collector must exclude %q", directory)
	}

	if strings.TrimSpace(repository.GeneratedMarker) == "" {
		return fmt.Errorf("collector generated-file marker must not be empty")
	}

	if repository.GeneratedProbeLimit <= 0 {
		return fmt.Errorf("collector generated-file probe limit must be positive")
	}

	return nil
}

/* ----------------------------------------- Collection ----------------------------------------- */

func CollectFiles(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
	extensions ...string,
) (paths []string, err error) {
	return collectFilteredFiles(repoRoot, repository, scope, func(path string) bool {
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

func CollectAllFiles(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (paths []string, err error) {
	return collectFilteredFiles(repoRoot, repository, scope, func(path string) bool {
		info, statErr := os.Stat(path)
		if statErr != nil {
			return false
		}

		return info.Mode().IsRegular()
	})
}

func collectFilteredFiles(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
	include func(path string) bool,
) (paths []string, err error) {
	roots := scopeRoots(repoRoot, repository, scope)

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

/* -------------------------------------------- Paths ------------------------------------------- */

func scopeRoots(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (roots []string) {
	return repository.ScanRoots(repoRoot, scope)
}

func RelativePath(repoRoot string, path string) (relative string) {
	relative, err := filepath.Rel(repoRoot, path)
	if err != nil {
		return filepath.Clean(path)
	}

	return filepath.ToSlash(relative)
}

func dedupePaths(values []string) (deduped []string) {
	seen := make(map[string]bool)
	deduped = make([]string, 0, len(values))

	for _, value := range values {
		if seen[value] {
			continue
		}

		seen[value] = true
		deduped = append(deduped, value)
	}

	return deduped
}

func isExcludedDirectory(repository profile.RepositoryConfig, name string) (excluded bool) {
	for _, current := range repository.GlobalExclusions {
		if current == name {
			return true
		}
	}

	return false
}

/* ------------------------------------- Generated Detection ------------------------------------ */

func isGeneratedFile(path string, repository profile.RepositoryConfig) (generated bool) {
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	buffer := make([]byte, repository.GeneratedProbeLimit)
	count, readErr := file.Read(buffer)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		if closeErr := file.Close(); closeErr != nil {
			return false
		}
		return false
	}

	if closeErr := file.Close(); closeErr != nil {
		return false
	}

	return hasGeneratedHeader(string(buffer[:count]), repository.GeneratedMarker)
}

func hasGeneratedHeader(contents string, marker string) (generated bool) {
	inspectedLines := 0
	inBlockComment := false

	for _, line := range strings.Split(contents, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		inspectedLines++
		if inspectedLines > generatedHeaderLineLimit {
			return false
		}

		if hasGeneratedHeaderLine(trimmed, marker, inBlockComment) {
			return true
		}

		if strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "<!--") {
			inBlockComment = true
		}

		if strings.Contains(trimmed, "*/") || strings.Contains(trimmed, "-->") {
			inBlockComment = false
		}
	}

	return false
}

func hasGeneratedHeaderLine(
	line string,
	marker string,
	inBlockComment bool,
) (generated bool) {
	comment := generatedCommentBody(line, inBlockComment)
	if comment == "" || !strings.Contains(comment, marker) {
		return false
	}

	return strings.Contains(strings.ToLower(comment), "generated")
}

func generatedCommentBody(line string, inBlockComment bool) (comment string) {
	if inBlockComment && strings.HasPrefix(line, "*") {
		comment = strings.TrimSpace(strings.TrimPrefix(line, "*"))
		comment = strings.TrimSuffix(comment, "*/")
		return strings.TrimSpace(comment)
	}

	if strings.HasPrefix(line, "/*") {
		comment = strings.TrimSpace(strings.TrimPrefix(line, "/*"))
		comment = strings.TrimSuffix(comment, "*/")
		return strings.TrimSpace(comment)
	}

	if strings.HasPrefix(line, "<!--") {
		comment = strings.TrimSpace(strings.TrimPrefix(line, "<!--"))
		comment = strings.TrimSuffix(comment, "-->")
		return strings.TrimSpace(comment)
	}

	for _, prefix := range generatedCommentPrefixes {
		if !strings.HasPrefix(line, prefix) {
			continue
		}

		comment = strings.TrimSpace(strings.TrimPrefix(line, prefix))
		return strings.TrimSpace(comment)
	}

	return ""
}
