package filewalk

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RelativePath returns path relative to root using forward slashes. It returns an error
// rather than an absolute or parent-escaping fallback: a path that cannot be contained under
// root is a containment failure that must surface. Repository-owned records resolve the
// error; DisplayPath provides a safe display fallback for diagnostics.
func RelativePath(root string, path string) (relative string, err error) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", fmt.Errorf("path %q is not relative to repository root", path)
	}

	cleaned := filepath.ToSlash(rel)
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", fmt.Errorf("path %q escapes the repository root", path)
	}

	return cleaned, nil
}

// DisplayPath returns path relative to root using forward slashes for diagnostic display.
// When the path cannot be relativised it falls back to the base name, so it never emits an
// absolute or escaping value. Callers that must reject escapes should use RelativePath.
func DisplayPath(root string, path string) (relative string) {
	relative, err := RelativePath(root, path)
	if err != nil {
		return filepath.ToSlash(filepath.Base(path))
	}

	return relative
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
