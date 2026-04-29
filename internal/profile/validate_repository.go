package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

func validateRepository(repository policy.RepositoryConfig) (err error) {
	if len(repository.RootMarkers) == 0 {
		return fmt.Errorf("repository.root_markers must not be empty")
	}

	if len(repository.Scopes) == 0 {
		return fmt.Errorf("repository.scopes must not be empty")
	}

	if repository.DefaultScope == "" {
		return fmt.Errorf("repository.default_scope must not be empty")
	}

	if repository.GeneratedMarker == "" {
		return fmt.Errorf("repository.generated_marker must not be empty")
	}

	if repository.GeneratedProbeLimit <= 0 {
		return fmt.Errorf("repository.generated_probe_limit must be positive")
	}

	for scope, roots := range repository.Scopes {
		if scope == "" {
			return fmt.Errorf("repository.scopes contains an empty scope")
		}

		if len(roots) == 0 {
			return fmt.Errorf("repository.scopes.%s must not be empty", scope)
		}
	}

	if !repository.ScopeExists(repository.DefaultScope) {
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

	if styleGuide.RequirementIDScheme != policy.RequirementIDSchemeSectionSlug {
		return fmt.Errorf(
			"unsupported styleguide.requirement_id_scheme %q",
			styleGuide.RequirementIDScheme,
		)
	}

	return nil
}

func validateImports(imports policy.ImportsConfig) (err error) {
	if imports.LocalPrefix == "" {
		return fmt.Errorf("imports.local_prefix must not be empty")
	}

	return nil
}
