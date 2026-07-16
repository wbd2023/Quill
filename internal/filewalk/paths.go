package filewalk

import "path/filepath"

// RelativePath returns path relative to repoRoot, using forward slashes. If the path cannot be
// made relative, it is returned cleaned with system separators.
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

func isExcludedDirectory(config WalkConfig, name string) (excluded bool) {
	for _, exclusion := range config.ExcludedDirectories {
		if exclusion == name {
			return true
		}
	}

	return false
}
