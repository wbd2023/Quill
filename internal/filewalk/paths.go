package filewalk

import (
	"path/filepath"
	"slices"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

func collectScopeRoots(
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

	slices.Sort(roots)
	return roots
}

// RelativePath relative path.
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

func isExcludedDirectory(repository policy.RepositoryConfig, name string) (excluded bool) {
	for _, exclusion := range repository.ExcludedDirectories {
		if exclusion == name {
			return true
		}
	}

	return false
}
