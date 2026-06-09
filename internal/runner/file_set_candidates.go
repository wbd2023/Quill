package runner

import (
	"os"
	"path/filepath"
	"slices"
	"sort"

	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

func collectFileSetCandidates(
	context Context,
	fileSet policy.FileSetConfig,
) (files []string, err error) {
	scopes := []style.Scope{context.Scope}
	if fileSetUsesScopedIncludes(fileSet) {
		scopes = findOverlappingScopes(
			context.Profile.Repository,
			context.Scope,
			fileSetIncludeScopes(fileSet),
		)
		if len(scopes) == 0 {
			return nil, nil
		}
	}

	files, err = filewalk.CollectFilesInScopes(
		context.RepoRoot,
		context.Profile.Repository,
		scopes,
		fileSet.Include.Extensions...,
	)
	if err != nil {
		return nil, err
	}

	files = append(files, explicitFileCandidates(context, fileSet, scopes)...)
	sort.Strings(files)
	return dedupeCandidatePaths(files), nil
}

func findOverlappingScopes(
	repository policy.RepositoryConfig,
	scope style.Scope,
	candidates []style.Scope,
) (scopes []style.Scope) {
	seen := make(map[style.Scope]bool, len(candidates))
	for _, candidate := range candidates {
		if seen[candidate] {
			continue
		}

		seen[candidate] = true
		if repository.HasScopeOverlap(scope, candidate) {
			scopes = append(scopes, candidate)
		}
	}

	slices.Sort(scopes)
	return scopes
}

func explicitFileCandidates(
	context Context,
	fileSet policy.FileSetConfig,
	scopes []style.Scope,
) (files []string) {
	for _, scope := range scopes {
		for _, file := range fileSet.Include.Files[scope] {
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
