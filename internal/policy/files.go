package policy

import "ciphera/tools/internal/contract"

type FileSetConfig struct {
	Name                 string
	Extensions           []string
	Files                map[contract.Scope][]string
	Prefixes             map[contract.Scope][]string
	ExcludedExtensions   []string
	ExcludedNames        []string
	ExcludedNamePrefixes []string
	SkipBinary           bool
}
