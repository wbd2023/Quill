package runner

import (
	"os"
	"path/filepath"
	"sort"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

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
