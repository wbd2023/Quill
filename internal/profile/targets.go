package profile

import (
	"fmt"

	"github.com/wbd2023/quill/internal/style"
)

// TargetConfigs defines the targets available to rule bindings.
type TargetConfigs []TargetConfig

// TargetConfig binds language-specific target settings to a repository scope.
type TargetConfig struct {
	Name             string
	Language         string
	Scope            style.Scope
	WorkingDirectory string
	FormatPaths      []string
	CheckPaths       []string
}

// Lookup returns the named target.
func (targets TargetConfigs) Lookup(name string) (target TargetConfig, found bool) {
	for _, candidate := range targets {
		if candidate.Name == name {
			return candidate, true
		}
	}

	return TargetConfig{}, false
}

/* ------------------------------------------- Targets ------------------------------------------ */

func validateTargets(
	repository RepositoryConfig,
	targets TargetConfigs,
) (err error) {
	seen := make(map[string]bool, len(targets))
	for _, target := range targets {
		if isBlank(target.Name) {
			return fmt.Errorf("target name must not be empty")
		}

		if seen[target.Name] {
			return fmt.Errorf("duplicate target %q", target.Name)
		}

		seen[target.Name] = true

		if isBlank(target.Language) {
			return fmt.Errorf("target %q must define language", target.Name)
		}

		if target.WorkingDirectory != "" && isBlank(target.WorkingDirectory) {
			return fmt.Errorf("target %q has a blank working_directory", target.Name)
		}

		if target.WorkingDirectory != "" {
			if err = validateRepoPath(target.WorkingDirectory); err != nil {
				return fmt.Errorf("target %q working_directory: %w", target.Name, err)
			}
		}

		if err = validateTargetPaths(
			target.Name,
			"format_paths",
			target.FormatPaths,
		); err != nil {
			return err
		}

		if err = validateTargetPaths(
			target.Name,
			"check_paths",
			target.CheckPaths,
		); err != nil {
			return err
		}

		if !repository.HasScope(target.Scope) {
			return fmt.Errorf(
				"target %q references unknown scope %q",
				target.Name,
				target.Scope,
			)
		}
	}

	return nil
}

func validateTargetPaths(
	name string,
	field string,
	paths []string,
) (err error) {
	seen := make(map[string]bool, len(paths))
	for _, path := range paths {
		if isBlank(path) {
			return fmt.Errorf("target %q %s contains an empty path", name, field)
		}

		if err = validateRepoPath(path); err != nil {
			return fmt.Errorf("target %q %s: %w", name, field, err)
		}

		if seen[path] {
			return fmt.Errorf(
				"target %q %s contains duplicate path %q",
				name,
				field,
				path,
			)
		}

		seen[path] = true
	}

	return nil
}
