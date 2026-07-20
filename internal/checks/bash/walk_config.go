package bash

import (
	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/policy"
)

func walkConfig(repository policy.RepositoryConfig) (config filewalk.WalkConfig) {
	return filewalk.WalkConfig{
		ExcludedDirectories: repository.ExcludedDirectories,
		GeneratedMarker:     repository.GeneratedMarker,
	}
}
