package profile

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

/* ------------------------------------------ File Sets ----------------------------------------- */

func validateFileSets(
	repository policy.RepositoryConfig,
	fileSets []policy.FileSetConfig,
) (err error) {
	seen := make(map[string]bool, len(fileSets))
	for _, fileSet := range fileSets {
		if isBlank(fileSet.Name) {
			return fmt.Errorf("file set name must not be empty")
		}

		if seen[fileSet.Name] {
			return fmt.Errorf("duplicate file set %q", fileSet.Name)
		}

		seen[fileSet.Name] = true

		if err = validateFileSetFilters(fileSet); err != nil {
			return err
		}

		if err = validateFileSetPaths(
			repository,
			fileSet.Name,
			"include.files",
			fileSet.Include.Files,
		); err != nil {
			return err
		}

		if err = validateFileSetPaths(
			repository,
			fileSet.Name,
			"include.paths",
			fileSet.Include.Paths,
		); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------------------- File Set Filters -------------------------------------- */

func validateFileSetFilters(fileSet policy.FileSetConfig) (err error) {
	if err = validateFileSetFilter(
		fileSet.Name,
		"include.extensions",
		fileSet.Include.Extensions,
	); err != nil {
		return err
	}

	if err = validateFileSetFilter(
		fileSet.Name,
		"exclude.extensions",
		fileSet.Exclude.Extensions,
	); err != nil {
		return err
	}

	return validateFileSetFilter(
		fileSet.Name,
		"exclude.files",
		fileSet.Exclude.Files,
	)
}

func validateFileSetFilter(
	fileSetName string,
	field string,
	values []string,
) (err error) {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if isBlank(value) {
			return fmt.Errorf("file set %q %s contains an empty value", fileSetName, field)
		}

		if seen[value] {
			return fmt.Errorf(
				"file set %q %s contains duplicate value %q",
				fileSetName,
				field,
				value,
			)
		}

		seen[value] = true
	}

	return nil
}

/* ---------------------------------------- Scoped Paths ---------------------------------------- */

func validateFileSetPaths(
	repository policy.RepositoryConfig,
	fileSetName string,
	field string,
	paths map[style.Scope][]string,
) (err error) {
	for scope, values := range paths {
		if isBlank(string(scope)) {
			return fmt.Errorf("file set %q %s contains an empty scope", fileSetName, field)
		}

		if !repository.HasScope(scope) {
			return fmt.Errorf("file set %q references unknown scope %q", fileSetName, scope)
		}

		if len(values) == 0 {
			return fmt.Errorf("file set %q %s.%s must not be empty", fileSetName, field, scope)
		}

		seen := make(map[string]bool, len(values))
		for _, value := range values {
			if isBlank(value) {
				return fmt.Errorf(
					"file set %q %s.%s contains an empty path",
					fileSetName,
					field,
					scope,
				)
			}

			if seen[value] {
				return fmt.Errorf(
					"file set %q %s.%s contains duplicate path %q",
					fileSetName,
					field,
					scope,
					value,
				)
			}

			seen[value] = true
		}
	}

	return nil
}
