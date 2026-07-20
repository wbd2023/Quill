package analysis

import (
	"path/filepath"
	"strings"

	"github.com/wbd2023/Quill/internal/policy"
)

// PathClassifier is path classifier.
type PathClassifier struct {
	repoRoot string
	paths    policy.PathRoles
}

// NewPathClassifier new path classifier.
func NewPathClassifier(
	repoRoot string,
	paths policy.PathRoles,
) (classifier PathClassifier) {
	return PathClassifier{
		repoRoot: filepath.Clean(repoRoot),
		paths:    paths,
	}
}

func (classifier PathClassifier) HasRole(path string, roleName string) (found bool) {
	relativePath := classifier.relativePath(path)
	for _, pattern := range classifier.paths.LookupPatterns(roleName) {
		if matchesRelativePathPattern(relativePath, pattern) {
			return true
		}
	}

	return false
}

func (classifier PathClassifier) MatchesImportPath(
	importPath string,
	roleName string,
) (found bool) {
	for _, pattern := range classifier.paths.LookupPatterns(roleName) {
		trimmedPattern := strings.TrimSuffix(filepath.ToSlash(pattern), "/")
		if trimmedPattern == "" {
			continue
		}

		if importPath == trimmedPattern || strings.HasSuffix(importPath, "/"+trimmedPattern) {
			return true
		}
	}

	return false
}

func (classifier PathClassifier) FirstPattern(roleName string) (pattern string) {
	patterns := classifier.paths.LookupPatterns(roleName)
	if len(patterns) == 0 {
		return ""
	}

	return patterns[0]
}

func (classifier PathClassifier) relativePath(path string) (relativePath string) {
	cleanedPath := filepath.Clean(path)
	if filepath.IsAbs(cleanedPath) {
		relativePath, err := filepath.Rel(classifier.repoRoot, cleanedPath)
		if err == nil {
			return filepath.ToSlash(filepath.Clean(relativePath))
		}
	}

	return filepath.ToSlash(filepath.Clean(cleanedPath))
}

func matchesRelativePathPattern(relativePath string, pattern string) (matches bool) {
	normalisedPath := filepath.ToSlash(filepath.Clean(relativePath))
	normalisedPattern := filepath.ToSlash(filepath.Clean(strings.TrimSuffix(pattern, "/")))
	if strings.HasSuffix(pattern, "/") {
		return normalisedPath == normalisedPattern ||
			strings.HasPrefix(normalisedPath, normalisedPattern+"/")
	}

	return normalisedPath == normalisedPattern
}
