package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

/* ----------------------------------------- Repository ----------------------------------------- */

func validateRepository(repository policy.RepositoryConfig) (err error) {
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

		if seen[marker] {
			return fmt.Errorf("repository.root_markers contains duplicate marker %q", marker)
		}

		seen[marker] = true
	}

	return nil
}

/* ------------------------------------------- Scopes ------------------------------------------- */

func validateRepositoryScopes(repository policy.RepositoryConfig) (err error) {
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

func validateGeneratedFilePolicy(repository policy.RepositoryConfig) (err error) {
	if isBlank(repository.GeneratedMarker) {
		return fmt.Errorf("repository.generated_marker must not be empty")
	}

	return nil
}
