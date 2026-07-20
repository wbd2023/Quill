package execution

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/policy"
)

/* ------------------------------------------ Inclusion ----------------------------------------- */

func fileSetCoversPath(
	context RunContext,
	fileSet policy.FileSetConfig,
	path string,
) (covered bool) {
	if len(fileSet.Include.Files) == 0 && len(fileSet.Include.Paths) == 0 {
		return true
	}

	foundScope := false
	for scope, explicitFiles := range fileSet.Include.Files {
		if !context.Profile.Repository.HasScopeOverlap(context.Scope, scope) {
			continue
		}

		foundScope = true
		if fileSetCoversRelativePath(
			context.RepoRoot,
			path,
			explicitFiles,
			fileSet.Include.Paths[scope],
		) {
			return true
		}
	}

	for scope, pathPrefixes := range fileSet.Include.Paths {
		if _, alreadyChecked := fileSet.Include.Files[scope]; alreadyChecked {
			continue
		}

		if !context.Profile.Repository.HasScopeOverlap(context.Scope, scope) {
			continue
		}

		foundScope = true
		if fileSetCoversRelativePath(context.RepoRoot, path, nil, pathPrefixes) {
			return true
		}
	}

	return !foundScope
}

func fileSetCoversRelativePath(
	repoRoot string,
	path string,
	explicitFiles []string,
	pathPrefixes []string,
) (covered bool) {
	if len(explicitFiles) == 0 && len(pathPrefixes) == 0 {
		return true
	}

	relativePath := filewalk.RelativePath(repoRoot, path)
	return slices.Contains(explicitFiles, relativePath) ||
		hasPrefixedPath(pathPrefixes, relativePath)
}

/* ------------------------------------------ Exclusion ----------------------------------------- */

func fileSetExcludesPath(fileSet policy.FileSetConfig, path string) (excluded bool) {
	base := filepath.Base(path)
	if hasMatchingSuffix(fileSet.Exclude.Extensions, path) {
		return true
	}

	return hasMatchingFilePattern(fileSet.Exclude.Files, base, path)
}

/* -------------------------------------- Matching Helpers -------------------------------------- */

func hasMatchingFilePattern(patterns []string, base string, path string) (found bool) {
	for _, pattern := range patterns {
		target := base
		if strings.Contains(pattern, string(filepath.Separator)) ||
			strings.Contains(pattern, "/") {
			target = filepath.ToSlash(path)
			pattern = filepath.ToSlash(pattern)
		}

		matched, err := filepath.Match(pattern, target)
		if err == nil && matched {
			return true
		}
	}

	return false
}

func hasMatchingSuffix(suffixes []string, path string) (found bool) {
	for _, suffix := range suffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}

	return false
}

func hasPrefixedPath(pathPrefixes []string, path string) (found bool) {
	for _, prefix := range pathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}
