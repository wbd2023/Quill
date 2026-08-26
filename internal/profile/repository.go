package profile

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wbd2023/quill/internal/style"
)

// RepositoryConfig defines repository scope roots, exclusions, and generated-file detection.
type RepositoryConfig struct {
	RootMarkers         []string
	ScopeRoots          map[style.Scope][]string
	DefaultScope        style.Scope
	ExcludedDirectories []string
	GeneratedMarker     string
}

// HasScope reports whether the repository defines the named scope.
func (r RepositoryConfig) HasScope(scope style.Scope) (found bool) {
	_, found = r.ScopeRoots[scope]
	return found
}

// ResolveScopeRoots returns the filesystem roots for a scope under the repository root.
func (r RepositoryConfig) ResolveScopeRoots(
	root string,
	scope style.Scope,
) (roots []string) {
	scopeRoots := r.ScopeRoots[scope]
	roots = make([]string, 0, len(scopeRoots))
	for _, scopeRoot := range scopeRoots {
		scopeRoot = cleanScopeRoot(scopeRoot)
		if scopeRoot == "." {
			roots = append(roots, root)
			continue
		}

		roots = append(roots, filepath.Join(root, scopeRoot))
	}

	return roots
}

// HasScopeOverlap reports whether two scopes cover any common root.
func (r RepositoryConfig) HasScopeOverlap(
	scope style.Scope,
	other style.Scope,
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

/* ----------------------------------------- Repository ----------------------------------------- */

func validateRepository(repository RepositoryConfig) (err error) {
	if err = validateRepositoryMarkers(repository.RootMarkers); err != nil {
		return err
	}

	if err = validateRepositoryScopes(repository); err != nil {
		return err
	}

	if err = validateRepositoryExclusions(repository.ExcludedDirectories); err != nil {
		return err
	}

	return validateGeneratedFilePolicy(repository)
}

/* ---------------------------------------- Root Markers ---------------------------------------- */

func validateRepositoryMarkers(markers []string) (err error) {
	if len(markers) == 0 {
		return fmt.Errorf("repository.root_markers must not be empty")
	}

	seen := make(map[string]bool, len(markers))
	for _, marker := range markers {
		if isBlank(marker) {
			return fmt.Errorf("repository.root_markers contains an empty marker")
		}

		if err = validateRepoPath(marker); err != nil {
			return fmt.Errorf("repository.root_markers: %w", err)
		}

		if seen[marker] {
			return fmt.Errorf("repository.root_markers contains duplicate marker %q", marker)
		}

		seen[marker] = true
	}

	return nil
}

/* ------------------------------------------- Scopes ------------------------------------------- */

func validateRepositoryScopes(repository RepositoryConfig) (err error) {
	if len(repository.ScopeRoots) == 0 {
		return fmt.Errorf("repository.scope_roots must not be empty")
	}

	if isBlank(string(repository.DefaultScope)) {
		return fmt.Errorf("repository.default_scope must not be empty")
	}

	for scope, roots := range repository.ScopeRoots {
		if err = validateScopeRoots(string(scope), roots); err != nil {
			return err
		}
	}

	if !repository.HasScope(repository.DefaultScope) {
		return fmt.Errorf(
			"repository.default_scope references unknown scope %q",
			repository.DefaultScope,
		)
	}

	return nil
}

func validateScopeRoots(scope string, roots []string) (err error) {
	if isBlank(scope) {
		return fmt.Errorf("repository.scope_roots contains an empty scope")
	}

	if len(roots) == 0 {
		return fmt.Errorf("repository.scope_roots.%s must not be empty", scope)
	}

	seen := make(map[string]bool, len(roots))
	for _, root := range roots {
		if isBlank(root) {
			return fmt.Errorf("repository.scope_roots.%s contains an empty root", scope)
		}

		if err = validateRepoPath(root); err != nil {
			return fmt.Errorf("repository.scope_roots.%s: %w", scope, err)
		}

		if seen[root] {
			return fmt.Errorf(
				"repository.scope_roots.%s contains duplicate root %q",
				scope,
				root,
			)
		}

		seen[root] = true
	}

	return nil
}

/* ----------------------------------------- Exclusions ----------------------------------------- */

func validateRepositoryExclusions(exclusions []string) (err error) {
	seen := make(map[string]bool, len(exclusions))
	for _, exclusion := range exclusions {
		if isBlank(exclusion) {
			return fmt.Errorf("repository.excluded_directories contains an empty exclusion")
		}

		if seen[exclusion] {
			return fmt.Errorf(
				"repository.excluded_directories contains duplicate exclusion %q",
				exclusion,
			)
		}

		seen[exclusion] = true
	}

	return nil
}

/* --------------------------------------- Generated Files -------------------------------------- */

func validateGeneratedFilePolicy(repository RepositoryConfig) (err error) {
	if isBlank(repository.GeneratedMarker) {
		return fmt.Errorf("repository.generated_marker must not be empty")
	}

	return nil
}
