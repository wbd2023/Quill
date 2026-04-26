package filewalk

import (
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func scopeRoots(
	repoRoot string,
	repository policy.RepositoryConfig,
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

func isExcludedDirectory(repository policy.RepositoryConfig, name string) (excluded bool) {
	for _, exclusion := range repository.GlobalExclusions {
		if exclusion == name {
			return true
		}
	}

	return false
}
