package profile

import (
	"fmt"

	"github.com/wbd2023/quill/internal/style"
)

// FileSets defines the configured file sets.
type FileSets []FileSetConfig

// FileSetConfig defines a named group of repository text files. Binary files are skipped by
// scanners before file set filters are applied.
type FileSetConfig struct {
	Name    string
	Include FileSetInclude
	Exclude FileSetExclude
}

// FileSetInclude defines files selected into a file set.
type FileSetInclude struct {
	Extensions []string
	Files      map[style.Scope][]string
	Paths      map[style.Scope][]string
}

// FileSetExclude defines files removed from a file set.
type FileSetExclude struct {
	Extensions []string
	Files      []string
}

// Lookup returns the named file set.
func (f FileSets) Lookup(name string) (fileSet FileSetConfig, found bool) {
	for _, candidate := range f {
		if candidate.Name == name {
			return candidate, true
		}
	}

	return FileSetConfig{}, false
}

// Clone returns a deep copy of fileSets.
func (f FileSets) Clone() (clone FileSets) {
	if f == nil {
		return nil
	}

	clone = make(FileSets, 0, len(f))
	for _, fileSet := range f {
		clone = append(clone, fileSet.Clone())
	}

	return clone
}

// Clone returns a deep copy of fileSet.
func (fileSet FileSetConfig) Clone() (clone FileSetConfig) {
	return FileSetConfig{
		Name: fileSet.Name,
		Include: FileSetInclude{
			Extensions: append([]string{}, fileSet.Include.Extensions...),
			Files:      cloneScopePaths(fileSet.Include.Files),
			Paths:      cloneScopePaths(fileSet.Include.Paths),
		},
		Exclude: FileSetExclude{
			Extensions: append([]string{}, fileSet.Exclude.Extensions...),
			Files:      append([]string{}, fileSet.Exclude.Files...),
		},
	}
}

func cloneScopePaths(source map[style.Scope][]string) (clone map[style.Scope][]string) {
	if source == nil {
		return nil
	}

	clone = make(map[style.Scope][]string, len(source))
	for scope, paths := range source {
		clone[scope] = append([]string{}, paths...)
	}

	return clone
}

/* ------------------------------------------ File Sets ----------------------------------------- */

func validateFileSets(
	repository RepositoryConfig,
	fileSets []FileSetConfig,
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

func validateFileSetFilters(fileSet FileSetConfig) (err error) {
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
	repository RepositoryConfig,
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

			if err = validateRepoPath(value); err != nil {
				return fmt.Errorf("file set %q %s.%s: %w", fileSetName, field, scope, err)
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
