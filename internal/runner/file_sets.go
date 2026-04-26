package runner

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const binaryProbeLimit = 4096

/* --------------------------------------- File Selection --------------------------------------- */

func CollectFileSetFiles(context Context, name string) (files []string, err error) {
	fileSet, found := context.Policy.FileSet(name)
	if !found {
		return nil, errUnknownFileSet(name)
	}

	candidates, err := collectFileSetCandidates(context, fileSet)
	if err != nil {
		return nil, err
	}

	for _, path := range candidates {
		if !fileSetCoversPath(context, fileSet, path) {
			continue
		}

		if fileSetExcludesPath(fileSet, path) {
			continue
		}

		if fileSet.SkipBinary && isBinaryFile(path) {
			continue
		}

		files = append(files, path)
	}

	return files, nil
}

func collectFileSetCandidates(
	context Context,
	fileSet policy.FileSetConfig,
) (files []string, err error) {
	scopes := []contract.Scope{context.Scope}
	if fileSetUsesScopedIncludes(fileSet) {
		scopes = context.Policy.Repository.OverlappingScopes(
			context.Scope,
			fileSetIncludeScopes(fileSet),
		)
		if len(scopes) == 0 {
			return nil, nil
		}
	}

	files, err = filewalk.CollectFilesInScopes(
		context.RepoRoot,
		context.Policy.Repository,
		scopes,
		fileSet.Extensions...,
	)
	if err != nil {
		return nil, err
	}

	files = append(files, explicitFileCandidates(context, fileSet, scopes)...)
	sort.Strings(files)
	return dedupeCandidatePaths(files), nil
}

func fileSetUsesScopedIncludes(
	fileSet policy.FileSetConfig,
) (usesFilters bool) {
	for _, files := range fileSet.Files {
		if len(files) > 0 {
			return true
		}
	}

	for _, prefixes := range fileSet.Prefixes {
		if len(prefixes) > 0 {
			return true
		}
	}

	return false
}

func fileSetIncludeScopes(fileSet policy.FileSetConfig) (scopes []contract.Scope) {
	seen := make(map[contract.Scope]bool)
	for scope, files := range fileSet.Files {
		if len(files) == 0 || seen[scope] {
			continue
		}

		seen[scope] = true
		scopes = append(scopes, scope)
	}

	for scope, prefixes := range fileSet.Prefixes {
		if len(prefixes) == 0 || seen[scope] {
			continue
		}

		seen[scope] = true
		scopes = append(scopes, scope)
	}

	return scopes
}

func explicitFileCandidates(
	context Context,
	fileSet policy.FileSetConfig,
	scopes []contract.Scope,
) (files []string) {
	for _, scope := range scopes {
		for _, file := range fileSet.Files[scope] {
			path := filepath.Join(context.RepoRoot, file)
			info, err := os.Stat(path)
			if err != nil || !info.Mode().IsRegular() {
				continue
			}

			files = append(files, filepath.Clean(path))
		}
	}

	return files
}

func dedupeCandidatePaths(paths []string) (deduped []string) {
	seen := make(map[string]bool, len(paths))
	for _, path := range paths {
		if seen[path] {
			continue
		}

		seen[path] = true
		deduped = append(deduped, path)
	}

	return deduped
}

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

/* -------------------------------------- Binary Detection -------------------------------------- */

func isBinaryFile(path string) (binary bool) {
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	buffer := make([]byte, binaryProbeLimit)
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

	return bytes.IndexByte(buffer[:count], 0) >= 0
}

func hasPrefixedPath(prefixes []string, path string) (found bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

func errUnknownFileSet(name string) (err error) {
	return errors.New("unknown file set " + name)
}
