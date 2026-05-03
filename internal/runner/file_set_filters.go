package runner

import (
	"path/filepath"
	"slices"
	"strings"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

func fileSetCoversPath(
	context Context,
	fileSet policy.FileSetConfig,
	path string,
) (covered bool) {
	if len(fileSet.ExplicitFiles) == 0 && len(fileSet.PathPrefixes) == 0 {
		return true
	}

	foundScope := false
	for scope, explicitFiles := range fileSet.ExplicitFiles {
		if !context.Policy.Repository.HasScopeOverlap(context.Scope, scope) {
			continue
		}

		foundScope = true
		if fileSetCoversRelativePath(
			context.RepoRoot,
			path,
			explicitFiles,
			fileSet.PathPrefixes[scope],
		) {
			return true
		}
	}

	for scope, pathPrefixes := range fileSet.PathPrefixes {
		if _, alreadyChecked := fileSet.ExplicitFiles[scope]; alreadyChecked {
			continue
		}

		if !context.Policy.Repository.HasScopeOverlap(context.Scope, scope) {
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

func fileSetExcludesPath(fileSet policy.FileSetConfig, path string) (excluded bool) {
	base := filepath.Base(path)
	if hasMatchingSuffix(fileSet.ExcludedExtensions, path) {
		return true
	}

	return slices.Contains(fileSet.ExcludedNames, base) ||
		hasPrefixedPath(fileSet.ExcludedNamePrefixes, base)
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
