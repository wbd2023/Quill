package execution

import (
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

/* ------------------------------------- File Set Collection ------------------------------------ */

func collectFileSetCandidates(
	context RunContext,
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

	roots := resolveScopeRoots(context.RepoRoot, context.Profile.Repository, scopes)
	files, err = filewalk.CollectFilesInRoots(
		filewalk.WalkConfig{
			ExcludedDirectories: context.Profile.Repository.ExcludedDirectories,
			GeneratedMarker:     context.Profile.Repository.GeneratedMarker,
		},
		roots,
		fileSet.Include.Extensions...,
	)
	if err != nil {
		return nil, err
	}

	files = append(files, explicitFileCandidates(context, fileSet, scopes)...)
	sort.Strings(files)
	return dedupeCandidatePaths(files), nil
}

/* -------------------------------------- Scope Resolution -------------------------------------- */

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
	context RunContext,
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

/* ---------------------------------------- Deduplication --------------------------------------- */

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

func resolveScopeRoots(
	repoRoot string,
	repository policy.RepositoryConfig,
	scopes []style.Scope,
) (roots []string) {
	seen := make(map[string]bool)
	for _, scope := range scopes {
		for _, root := range repository.ResolveScopeRoots(repoRoot, scope) {
			clean := filepath.Clean(root)
			if seen[clean] {
				continue
			}

			seen[clean] = true
			roots = append(roots, clean)
		}
	}

	sort.Strings(roots)
	return roots
}
