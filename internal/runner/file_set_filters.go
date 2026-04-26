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
	if len(fileSet.Files) == 0 && len(fileSet.Prefixes) == 0 {
		return true
	}

	foundScope := false
	for scope, files := range fileSet.Files {
		if !context.Policy.Repository.ScopesOverlap(context.Scope, scope) {
			continue
		}

		foundScope = true
		if fileSetCoversRelativePath(context.RepoRoot, path, files, fileSet.Prefixes[scope]) {
			return true
		}
	}

	for scope, prefixes := range fileSet.Prefixes {
		if _, alreadyChecked := fileSet.Files[scope]; alreadyChecked {
			continue
		}

		if !context.Policy.Repository.ScopesOverlap(context.Scope, scope) {
			continue
		}

		foundScope = true
		if fileSetCoversRelativePath(context.RepoRoot, path, nil, prefixes) {
			return true
		}
	}

	return !foundScope
}

func fileSetCoversRelativePath(
	repoRoot string,
	path string,
	files []string,
	prefixes []string,
) (covered bool) {
	if len(files) == 0 && len(prefixes) == 0 {
		return true
	}

	relativePath := filewalk.RelativePath(repoRoot, path)
	return slices.Contains(files, relativePath) || hasPrefixedPath(prefixes, relativePath)
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

func hasPrefixedPath(prefixes []string, path string) (found bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}
