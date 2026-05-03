package policy

import "ciphera/tools/internal/contract"

// FileSets defines the configured file sets.
type FileSets []FileSetConfig

// FileSetConfig defines a named group of repository files.
type FileSetConfig struct {
	Name                 string
	Extensions           []string
	ExplicitFiles        map[contract.Scope][]string
	PathPrefixes         map[contract.Scope][]string
	ExcludedExtensions   []string
	ExcludedNames        []string
	ExcludedNamePrefixes []string
	SkipBinary           bool
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
