package runner

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/profile"
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
		if !fileSetCoversPath(context.RepoRoot, context.Scope, fileSet, path) {
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
	fileSet profile.FileSetConfig,
) (files []string, err error) {
	scope := context.Scope
	if fileSetUsesExplicitScope(context.Scope, fileSet) {
		scope = contract.ScopeAll
	}

	if len(fileSet.Extensions) == 0 {
		return filewalk.CollectAllFiles(context.RepoRoot, context.Policy.Repository, scope)
	}

	return filewalk.CollectFiles(
		context.RepoRoot,
		context.Policy.Repository,
		scope,
		fileSet.Extensions...,
	)
}

func fileSetUsesExplicitScope(
	scope contract.Scope,
	fileSet profile.FileSetConfig,
) (usesFilters bool) {
	switch scope {
	case contract.ScopeApp:
		return len(fileSet.AppFiles) > 0 || len(fileSet.AppPrefixes) > 0
	case contract.ScopeTools:
		return len(fileSet.ToolsFiles) > 0 || len(fileSet.ToolsPrefixes) > 0
	default:
		return false
	}
}

func fileSetCoversPath(
	repoRoot string,
	scope contract.Scope,
	fileSet profile.FileSetConfig,
	path string,
) (covered bool) {
	switch scope {
	case contract.ScopeApp:
		return fileSetCoversRelativePath(
			repoRoot,
			path,
			fileSet.AppFiles,
			fileSet.AppPrefixes,
		)

	case contract.ScopeTools:
		return fileSetCoversRelativePath(
			repoRoot,
			path,
			fileSet.ToolsFiles,
			fileSet.ToolsPrefixes,
		)

	default:
		return true
	}
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
	return containsString(files, relativePath) || hasPrefixedPath(prefixes, relativePath)
}

func fileSetExcludesPath(fileSet profile.FileSetConfig, path string) (excluded bool) {
	base := filepath.Base(path)
	if hasMatchingSuffix(fileSet.ExcludedExtensions, path) {
		return true
	}

	return containsString(fileSet.ExcludedNames, base) ||
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

func containsString(values []string, expected string) (found bool) {
	for _, value := range values {
		if value == expected {
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

func errUnknownFileSet(name string) (err error) {
	return errors.New("unknown file set " + name)
}
