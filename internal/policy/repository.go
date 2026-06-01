package policy

import (
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
)

// RepositoryConfig defines repository scope roots, exclusions, and generated-file detection.
type RepositoryConfig struct {
	RootMarkers         []string
	ScopeRoots          map[contract.Scope][]string
	DefaultScope        contract.Scope
	ExcludedDirectories []string
	GeneratedMarker     string
}

// HasScope reports whether the repository defines the named scope.
func (r RepositoryConfig) HasScope(scope contract.Scope) (found bool) {
	_, found = r.ScopeRoots[scope]
	return found
}

// ResolveScopeRoots returns the filesystem roots for a scope under the repository root.
func (r RepositoryConfig) ResolveScopeRoots(
	repositoryRoot string,
	scope contract.Scope,
) (roots []string) {
	scopeRoots := r.ScopeRoots[scope]
	roots = make([]string, 0, len(scopeRoots))
	for _, scopeRoot := range scopeRoots {
		scopeRoot = cleanScopeRoot(scopeRoot)
		if scopeRoot == "." {
			roots = append(roots, repositoryRoot)
			continue
		}

		roots = append(roots, filepath.Join(repositoryRoot, scopeRoot))
	}

	return roots
}

// HasScopeOverlap reports whether two scopes cover any common root.
func (r RepositoryConfig) HasScopeOverlap(
	scope contract.Scope,
	other contract.Scope,
) (overlap bool) {
	scopeRoots, otherRoots := r.ScopeRoots[scope], r.ScopeRoots[other]
	for _, scopeRoot := range scopeRoots {
		for _, otherRoot := range otherRoots {
			if hasRootOverlap(scopeRoot, otherRoot) {
				return true
			}
		}
	}

	return false
}

func hasRootOverlap(left string, right string) (overlap bool) {
	left, right = cleanScopeRoot(left), cleanScopeRoot(right)
	if left == "." || right == "." {
		return true
	}

	return left == right || strings.HasPrefix(left, right+"/") || strings.HasPrefix(right, left+"/")
}

func cleanScopeRoot(root string) (cleaned string) {
	root = strings.TrimSpace(root)
	if root == "" {
		return "."
	}

	return filepath.ToSlash(filepath.Clean(root))
}
