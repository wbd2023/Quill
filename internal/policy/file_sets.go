package policy

import "ciphera/tools/internal/style"

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
