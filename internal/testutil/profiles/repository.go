package profiles

import (
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

// RepositoryConfig returns a deterministic repository collector policy fixture.
func RepositoryConfig() (repository profile.RepositoryConfig) {
	return profile.RepositoryConfig{
		RootMarkers: []string{
			"STYLE.md",
			"quill.toml",
		},
		ScopeRoots: map[style.Scope][]string{
			"all": {"."},
		},
		DefaultScope: style.Scope("all"),
		ExcludedDirectories: []string{
			".cache",
			".git",
			".toolchain",
			"bin",
			"testdata",
			"third_party",
			"vendor",
		},
		GeneratedMarker: "DO NOT EDIT.",
	}
}
