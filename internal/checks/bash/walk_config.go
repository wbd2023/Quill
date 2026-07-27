package bash

import (
	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/policy"
)

func walkConfig(repository policy.RepositoryConfig) (config filewalk.WalkConfig) {
	return filewalk.WalkConfig{
		ExcludedDirectories: repository.ExcludedDirectories,
		GeneratedMarker:     repository.GeneratedMarker,
	}
}
