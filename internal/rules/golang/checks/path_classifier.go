package checks

import (
	"path/filepath"
	"strings"

	"ciphera/tools/internal/policy"
)

type PathClassifier struct {
	repoRoot string
	paths    policy.PathClassSet
}

func NewPathClassifier(
	repoRoot string,
	paths policy.PathClassSet,
) (classifier PathClassifier) {
	return PathClassifier{
		repoRoot: filepath.Clean(repoRoot),
		paths:    paths,
	}
}

func (classifier PathClassifier) HasClass(path string, className string) (found bool) {
	relativePath := classifier.relativePath(path)
	for _, pattern := range classifier.paths.Patterns(className) {
		if matchesRelativePathPattern(relativePath, pattern) {
			return true
		}
	}

	return false
}

func (classifier PathClassifier) MatchesImportPath(
	importPath string,
	className string,
) (found bool) {
	for _, pattern := range classifier.paths.Patterns(className) {
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

func (classifier PathClassifier) FirstPattern(className string) (pattern string) {
	patterns := classifier.paths.Patterns(className)
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
