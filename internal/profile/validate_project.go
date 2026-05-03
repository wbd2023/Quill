package profile

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/requirementid"
)

func validateRepository(repository policy.RepositoryConfig) (err error) {
	if len(repository.RootMarkers) == 0 {
		return fmt.Errorf("repository.root_markers must not be empty")
	}

	for _, marker := range repository.RootMarkers {
		if strings.TrimSpace(marker) == "" {
			return fmt.Errorf("repository.root_markers contains an empty marker")
		}
	}

	if len(repository.ScopeRoots) == 0 {
		return fmt.Errorf("repository.scope_roots must not be empty")
	}

	if repository.DefaultScope == "" {
		return fmt.Errorf("repository.default_scope must not be empty")
	}

	if repository.GeneratedMarker == "" {
		return fmt.Errorf("repository.generated_marker must not be empty")
	}

	if repository.GeneratedProbeBytes <= 0 {
		return fmt.Errorf("repository.generated_probe_bytes must be positive")
	}

	for scope, roots := range repository.ScopeRoots {
		if scope == "" {
			return fmt.Errorf("repository.scope_roots contains an empty scope")
		}

		if len(roots) == 0 {
			return fmt.Errorf("repository.scope_roots.%s must not be empty", scope)
		}

		for _, root := range roots {
			if strings.TrimSpace(root) == "" {
				return fmt.Errorf("repository.scope_roots.%s contains an empty root", scope)
			}
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

func validateStyleGuide(styleGuide policy.StyleGuideConfig) (err error) {
	if styleGuide.Path == "" {
		return fmt.Errorf("styleguide.path must not be empty")
	}

	if styleGuide.RequirementIDScheme == "" {
		return fmt.Errorf("styleguide.requirement_id_scheme must not be empty")
	}

	if styleGuide.RequirementIDScheme != requirementid.SectionSlug {
		return fmt.Errorf(
			"unsupported styleguide.requirement_id_scheme %q",
			styleGuide.RequirementIDScheme,
		)
	}

	return nil
}
