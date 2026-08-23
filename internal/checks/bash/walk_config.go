package bash

import (
	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/profile"
)

func walkConfig(repository profile.RepositoryConfig) (config filewalk.WalkConfig) {
	return filewalk.WalkConfig{
		ExcludedDirectories: repository.ExcludedDirectories,
		GeneratedMarker:     repository.GeneratedMarker,
	}
}
