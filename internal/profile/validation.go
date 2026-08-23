package profile

import (
	"fmt"
	"strings"
)

// Validate checks config for supported schema version and internal consistency.
func Validate(config Profile) (err error) {
	if config.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported style profile version %d", config.SchemaVersion)
	}

	if err = validateRepository(config.Repository); err != nil {
		return err
	}

	if err = validateStyleGuide(config.StyleGuide); err != nil {
		return err
	}

	if err = validatePathRoles(config.PathRoles); err != nil {
		return err
	}

	if err = validateFileSets(config.Repository, config.FileSets); err != nil {
		return err
	}

	if err = validateTargets(config.Repository, config.Targets); err != nil {
		return err
	}

	if err = validateTools(config.Tools); err != nil {
		return err
	}

	if err = validateEnabledPacks(config.EnabledPacks); err != nil {
		return err
	}
	if err = validatePackPolicies(config.EnabledPacks, config.PackPolicies); err != nil {
		return err
	}

	if err = validatePackSources(config.PackSources); err != nil {
		return err
	}

	if err = validateRules(
		config.Repository,
		config.Rules,
	); err != nil {
		return err
	}

	return nil
}

func isBlank(value string) (blank bool) {
	return strings.TrimSpace(value) == ""
}

func validatePackSources(sources []PackSource) (err error) {
	seen := make(map[string]bool, len(sources))
	for _, source := range sources {
		if source.Path == "" {
			return fmt.Errorf("pack_sources.path must not be empty")
		}

		if err = validateRepoPath(source.Path); err != nil {
			return fmt.Errorf("pack_sources.path %q: %w", source.Path, err)
		}

		if seen[source.Path] {
			return fmt.Errorf("pack_sources.path %q is declared more than once", source.Path)
		}
		seen[source.Path] = true
	}

	return nil
}
