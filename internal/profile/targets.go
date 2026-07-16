package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

/* ------------------------------------------- Targets ------------------------------------------ */

func validateTargets(
	repository policy.RepositoryConfig,
	targets policy.TargetConfigs,
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
