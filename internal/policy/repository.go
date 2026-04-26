package policy

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"ciphera/tools/internal/contract"
)

/* -------------------------------------------- Types ------------------------------------------- */

type RepositoryConfig struct {
	RootMarkers         []string
	DefaultScope        contract.Scope
	Scopes              map[contract.Scope][]string
	GlobalExclusions    []string
	GeneratedMarker     string
	GeneratedProbeLimit int
}

/* ------------------------------------------- Scopes ------------------------------------------- */

func (repository RepositoryConfig) ScopeExists(scope contract.Scope) (found bool) {
	_, found = repository.Scopes[scope]
	return found
}

func (repository RepositoryConfig) ScanRoots(
	repoRoot string,
	scope contract.Scope,
) (roots []string) {
	return joinPaths(repoRoot, repository.Scopes[scope])
}

func (repository RepositoryConfig) ScanRootsForScopes(
	repoRoot string,
	scopes []contract.Scope,
) (roots []string) {
	seen := make(map[string]bool)
	for _, scope := range scopes {
		for _, root := range repository.ScanRoots(repoRoot, scope) {
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

func (repository RepositoryConfig) OverlappingScopes(
	scope contract.Scope,
	candidates []contract.Scope,
) (scopes []contract.Scope) {
	seen := make(map[contract.Scope]bool, len(candidates))
	for _, candidate := range candidates {
		if seen[candidate] {
			continue
		}

		seen[candidate] = true
		if repository.ScopesOverlap(scope, candidate) {
			scopes = append(scopes, candidate)
		}
	}

	slices.Sort(scopes)
	return scopes
}

func (repository RepositoryConfig) ScopesOverlap(
	left contract.Scope,
	right contract.Scope,
) (overlap bool) {
	leftRoots := repository.Scopes[left]
	rightRoots := repository.Scopes[right]
	for _, leftRoot := range leftRoots {
		for _, rightRoot := range rightRoots {
			if scopeRootsOverlap(leftRoot, rightRoot) {
				return true
			}
		}
	}

	return false
}

/* ----------------------------------------- Validation ----------------------------------------- */

func (repository RepositoryConfig) ValidateRoot(repoRoot string) (err error) {
	for _, marker := range repository.RootMarkers {
		if marker == "" {
			continue
		}

		if _, statErr := os.Stat(filepath.Join(repoRoot, marker)); statErr == nil {
			continue
		} else {
			return statErr
		}
	}

	return nil
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func joinPaths(repoRoot string, values []string) (paths []string) {
	paths = make([]string, 0, len(values))

	for _, value := range values {
		if value == "." {
			paths = append(paths, repoRoot)
			continue
		}

		paths = append(paths, filepath.Join(repoRoot, value))
	}

	return paths
}

func scopeRootsOverlap(left string, right string) (overlap bool) {
	left = normaliseScopeRoot(left)
	right = normaliseScopeRoot(right)
	if left == "." || right == "." {
		return true
	}

	return left == right ||
		strings.HasPrefix(left, right+"/") ||
		strings.HasPrefix(right, left+"/")
}

func normaliseScopeRoot(root string) (normalised string) {
	normalised = filepath.ToSlash(filepath.Clean(strings.TrimSpace(root)))
	if normalised == "" {
		return "."
	}

	return normalised
}
